package chain

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type msgPublisherCmdMsg struct {
	act *ActMsgHandler
	trx *TrxMsgHandler
	blk *BlkMsgHandler
}

// messagePublisher chain message publisher
type messagePublisher struct {
	logger *zap.Logger

	actHandlers []*ActMsgHandler
	trxHandlers []*TrxMsgHandler
	blkHandlers []*BlkMsgHandler

	// msg channel
	blkChan chan *SignedBlock
	cmdChan chan msgPublisherCmdMsg

	wg sync.WaitGroup
}

func newMessagePublisher(ctx context.Context, logger *zap.Logger) *messagePublisher {
	res := &messagePublisher{
		logger:      logger,
		actHandlers: make([]*ActMsgHandler, 0, 16),
		trxHandlers: make([]*TrxMsgHandler, 0, 16),
		blkHandlers: make([]*BlkMsgHandler, 0, 16),
		blkChan:     make(chan *SignedBlock, 2048),
		cmdChan:     make(chan msgPublisherCmdMsg, 64),
	}

	res.start(ctx)

	return res
}

func (m *messagePublisher) start(ctx context.Context) {
	m.wg.Add(1)
	go m.loop(ctx)
}

func (m *messagePublisher) loop(ctx context.Context) {
	defer m.wg.Done()
	for {
		select {
		case cmdMsg := <-m.cmdChan:
			m.onCmdMsg(&cmdMsg)

		case blk := <-m.blkChan:
			m.onCommitedBlock(blk)

		case <-ctx.Done():
			m.logger.Info("message publisher closed")
			return
		default:
		}
	}
}

func (m *messagePublisher) onCmdMsg(msg *msgPublisherCmdMsg) {
	if msg.act != nil {
		m.actHandlers = append(m.actHandlers, msg.act)
	}

	if msg.trx != nil {
		m.trxHandlers = append(m.trxHandlers, msg.trx)
	}

	if msg.blk != nil {
		m.blkHandlers = append(m.blkHandlers, msg.blk)
	}
}

func (m *messagePublisher) onCommitedBlock(blk *SignedBlock) {
	n := blk.BlockNumber()
	if n%1000 == 0 {
		m.logger.Debug("on committed block", zap.Uint32("num", blk.BlockNumber()))
	}
}

// AppendActHandler append a action handler
func (m *messagePublisher) AppendActHandler(h ActMsgHandler) error {
	m.cmdChan <- msgPublisherCmdMsg{
		act: &h,
	}
	return nil
}

// AppendTrxHandler append a trx handler
func (m *messagePublisher) AppendTrxHandler(h TrxMsgHandler) error {
	m.cmdChan <- msgPublisherCmdMsg{
		trx: &h,
	}
	return nil
}

// AppendBlockHandler append a action handler
func (m *messagePublisher) AppendBlockHandler(h BlkMsgHandler) error {
	m.cmdChan <- msgPublisherCmdMsg{
		blk: &h,
	}
	return nil
}

// OnCommittedBlock committed block to publisher
func (m *messagePublisher) OnCommittedBlock(b *SignedBlock) error {
	m.blkChan <- b
	return nil
}

// Wait wait to exit
func (m *messagePublisher) Wait() {
	m.wg.Wait()
}
