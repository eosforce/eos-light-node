package p2p

import (
	"context"
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
	isClosing bool
	clientsWg sync.WaitGroup
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

func (p *p2pClientImp) Start(ctx context.Context) error {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case ev, ok := <-p.msgChan:
				if !ok {
					p.logger.Warn("p2p peers msg chan closed")
					return
				}
				p.logger.Debug("handler msg")
				p.handlerImp.handleImp(ev)

			case <-ctx.Done():
				if !p.isClosing {
					p.isClosing = true
					p.logger.Info("p2p client imp start stop")
					p.closeConnection()
					p.logger.Debug("wait clients close")
					//p.clientsWg.Wait()
					//p.logger.Info("close msg channel")
					//close(p.msgChan)
					return
				}
			default:
			}
		}
	}()

	for idx, client := range p.clients {
		p.createClient(ctx, idx, client)
	}

	return nil
}

func (p *p2pClientImp) IsClosed() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.hasClosed
}

func (p *p2pClientImp) createClient(ctx context.Context, idx int, client p2pClientInterface) {
	p.clientsWg.Add(1)
	go func() {
		defer p.clientsWg.Done()
		for {
			p.logger.Info("create connect", zap.Int("client", idx))
			err := client.Start()

			// check when after close client
			if p.IsClosed() {
				p.logger.Info("client closed", zap.Int("client", idx))
				return
			}

			if err != nil {
				p.logger.Error("client start err", zap.Int("client", idx), zap.Error(err))
			}

			p.logger.Info("client wait to reconnect", zap.Int("client", idx))

			// wait to reconnect
			time.Sleep(3 * time.Second)

			// check when after sleep
			if p.IsClosed() {
				p.logger.Info("client closed", zap.Int("client", idx))
				return
			}
		}
	}()
}

func (p *p2pClientImp) closeConnection() error {
	p.mutex.Lock()
	p.hasClosed = true
	p.mutex.Unlock()

	for idx, client := range p.clients {
		err := client.CloseConnection()
		if err != nil {
			p.logger.Error("client close err", zap.Int("client", idx), zap.Error(err))
		}
		p.logger.Info("client close", zap.Int("client", idx))
	}

	return nil
}

func (p *p2pClientImp) onMsg(envelope p2pClientMsg) {
	p.msgChan <- envelope
}

func (p *p2pClientImp) RegHandler(handler p2pHandler) {
	p.handlers = append(p.handlers, handler)
}

func (p *p2pClientImp) SetReadTimeout(readTimeout time.Duration) {
	for _, peer := range p.clients {
		peer.SetReadTimeout(readTimeout)
	}
}
