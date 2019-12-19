package p2p

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/fanyang1988/eos-light-node/core/chain"
	"github.com/fanyang1988/eos-p2p/p2p"
	"go.uber.org/zap"
)

// Client a manager for peers to diff p2p node
type Client struct {
	*p2pClientImp
}

// NewP2PClient new p2p peers from cfg
func NewP2PClient(ctx context.Context, name string, chainID string, startBlock uint32, peers []string, logger *zap.Logger) *Client {
	p := &Client{
		&p2pClientImp{},
	}

	p.init(name, chainID, peers, logger)
	p.setHandlerImp(p)

	peersCfg := make([]*p2p.PeerCfg, 0, len(peers))
	for _, p := range peers {
		peersCfg = append(peersCfg, &p2p.PeerCfg{
			Address: p,
		})
	}

	//Logger.Info("P2P Client ", zap.String("peer", *peer), zap.String("chainID", *chainID))
	client, err := p2p.NewClient(
		ctx,
		chainID,
		peersCfg,
		p2p.WithNeedSync(startBlock),
		p2p.WithHandler(p2p.StringLoggerHandler),
		p2p.WithHandler(p),
	)

	if err != nil {
		panic(err)
	}

	p.client = client

	return p
}

// Wait wait client to stop
func (p *Client) Wait() {
	p.client.Wait()
	p.wg.Wait()
}

// Name imp handler
func (p *Client) Name() string {
	return "node-p2p"
}

func (p *Client) handleImp(m p2pClientMsg) {
	peer := m.peer
	pkg, ok := m.msg.(*chain.Packet)
	if !ok {
		p.logger.Error("packet type err")
		return
	}

	for _, h := range p.handlers {
		func(hh p2pHandler) {
			defer func() {
				if err := recover(); err != nil {
					p.logger.Error("handler process ev panic",
						zap.String("err", fmt.Sprintf("err:%s", err)),
						zap.String("stack", string(debug.Stack())))
				}
			}()

			var err error
			switch pkg.Type {
			case chain.GoAwayMessageType:
				m, ok := pkg.P2PMessage.(*chain.GoAwayMessage)
				if !ok {
					p.logger.Error("msg type err by go away")
					return
				}
				p.logger.Info("peer goaway",
					zap.String("peer", peer),
					zap.String("reason", m.Reason.String()),
					zap.String("nodeId", m.NodeID.String()))
				err = hh.OnGoAway(peer, uint8(m.Reason), m.NodeID)
			case chain.SignedBlockType:
				m, ok := pkg.P2PMessage.(*chain.SignedBlock)
				if !ok {
					p.logger.Error("msg type err by go away")
					return
				}
				p.logger.Debug("on signed block",
					zap.String("peer", peer),
					zap.String("block", m.String()))
				if err == nil {
					err = hh.OnBlock(peer, m)
				} else {
					p.logger.Error("handle msg err", zap.Error(err))
				}
			}

			if err != nil {
				p.logger.Error("handle msg err", zap.Error(err))
			}

		}(h)
	}
}

// Handle handler for p2p clients
func (p *Client) Handle(envelope *p2p.Envelope) {
	p.onMsg(p2pClientMsg{
		peer: envelope.Sender.Address,
		msg:  envelope.Packet,
	})
}
