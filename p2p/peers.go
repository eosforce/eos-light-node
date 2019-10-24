package p2p

import (
	"encoding/hex"
	"fmt"
	"runtime/debug"

	eos "github.com/eosforce/goeosforce"
	"github.com/eosforce/goeosforce/p2p"
	"go.uber.org/zap"
)

// P2PClient a manager for peers to diff p2p node
type P2PClient struct {
	*p2pClientImp
}

// NewP2PClient new p2p peers from cfg
func NewP2PClient(name string, chainID string, startBlock uint32, peers []string, logger *zap.Logger) *P2PClient {
	p := &P2PClient{
		&p2pClientImp{},
	}

	p.init(name, chainID, peers, logger)
	p.setHandlerImp(p)

	cID, err := hex.DecodeString(chainID)
	if err != nil {
		p.logger.Error("decode chain id err", zap.Error(err))
		panic(err)
	}

	for idx, peer := range peers {
		p.logger.Debug("new peer client", zap.Int("idx", idx), zap.String("peer", peer))
		client := p2p.NewClient(
			p2p.NewOutgoingPeer(peer, fmt.Sprintf("%s-%02d", name, idx), &p2p.HandshakeInfo{
				ChainID:      cID,
				HeadBlockNum: startBlock,
			}),
			true,
		)
		client.RegisterHandler(p)
		p.clients = append(p.clients, client)
	}

	return p
}

func (p *P2PClient) handleImp(m p2pClientMsg) {
	peer := m.peer
	pkg, ok := m.msg.(*eos.Packet)
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
			case eos.GoAwayMessageType:
				m, ok := pkg.P2PMessage.(*eos.GoAwayMessage)
				if !ok {
					p.logger.Error("msg type err by go away")
					return
				}
				p.logger.Info("peer goaway",
					zap.String("peer", peer),
					zap.String("reason", m.Reason.String()),
					zap.String("nodeId", m.NodeID.String()))
				err = hh.OnGoAway(peer, uint8(m.Reason), m.NodeID)
			case eos.SignedBlockType:
				m, ok := pkg.P2PMessage.(*eos.SignedBlock)
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
func (p *P2PClient) Handle(envelope *p2p.Envelope) {
	p.onMsg(p2pClientMsg{
		peer: envelope.Sender.Address,
		msg:  envelope.Packet,
	})
}
