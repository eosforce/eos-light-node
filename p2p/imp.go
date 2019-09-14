package p2p

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type p2pClientInterface interface {
	Start() error
	CloseConnection() error
	SetReadTimeout(readTimeout time.Duration)
}

type p2pClientMsg struct {
	peer   string
	msg    interface{}
	isStop bool
}

type p2pHandlerInterface interface {
	handleImp(msg p2pClientMsg)
}

// p2pForceioClient a manager for peers to diff p2p node
type p2pClientImp struct {
	name      string
	clients   []p2pClientInterface
	handlers  []p2pHandler
	msgChan   chan p2pClientMsg
	wg        sync.WaitGroup
	chanWg    sync.WaitGroup
	hasClosed bool
	mutex     sync.RWMutex
	logger    *zap.Logger

	handlerImp p2pHandlerInterface
}

func (p *p2pClientImp) init(name string, chainID string, peers []string, logger *zap.Logger) {
	p.name = name
	p.clients = make([]p2pClientInterface, 0, len(peers))
	p.handlers = make([]p2pHandler, 0, 8)
	p.msgChan = make(chan p2pClientMsg, 4096)
	p.logger = logger
}

func (p *p2pClientImp) setHandlerImp(h p2pHandlerInterface) {
	p.handlerImp = h
}

func (p *p2pClientImp) Start() error {
	p.chanWg.Add(1)
	go func() {
		defer p.chanWg.Done()
		for {
			isStop := p.Loop()
			if isStop {
				p.logger.Info("p2p peers stop")
				return
			}
		}
	}()

	for idx, client := range p.clients {
		p.createClient(idx, client)
	}

	return nil
}

func (p *p2pClientImp) IsClosed() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.hasClosed
}

func (p *p2pClientImp) createClient(idx int, client p2pClientInterface) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			p.logger.Info("create connect", zap.Int("client", idx))
			err := client.Start()

			// check when after close client
			if p.IsClosed() {
				return
			}

			if err != nil {
				p.logger.Error("client err", zap.Int("client", idx), zap.Error(err))
			}

			time.Sleep(3 * time.Second)

			// check when after sleep
			if p.IsClosed() {
				return
			}
		}
	}()
}

func (p *p2pClientImp) CloseConnection() error {
	p.logger.Warn("start close")

	p.mutex.Lock()
	p.hasClosed = true
	p.mutex.Unlock()

	for idx, client := range p.clients {
		go func(i int, cli p2pClientInterface) {
			err := cli.CloseConnection()
			if err != nil {
				p.logger.Error("client close err", zap.Int("client", i), zap.Error(err))
			}
			p.logger.Info("client close", zap.Int("client", i))
		}(idx, client)
	}
	p.wg.Wait()
	p.msgChan <- p2pClientMsg{
		isStop: true,
	}
	close(p.msgChan)
	p.chanWg.Wait()

	return nil
}

func (p *p2pClientImp) Loop() bool {
	ev, ok := <-p.msgChan
	if ev.isStop {
		return true
	}

	if !ok {
		p.logger.Warn("p2p peers msg chan closed")
		return true
	}

	p.handlerImp.handleImp(ev)

	return false
}

func (p *p2pClientImp) onMsg(envelope p2pClientMsg) {
	p.msgChan <- envelope
}

func (p *p2pClientImp) RegHandler(handler p2pHandler) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	p.handlers = append(p.handlers, handler)
}

func (p *p2pClientImp) SetReadTimeout(readTimeout time.Duration) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, peer := range p.clients {
		peer.SetReadTimeout(readTimeout)
	}
}
