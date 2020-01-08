package main

import (
	"github.com/eosforce/eos-light-node/core/chain"
	"github.com/eosforce/eos-p2p/p2p"
)

type chainP2PHandler struct {
	chain *chain.Chain
}

func (h *chainP2PHandler) Handle(envelope *p2p.Envelope) {
	signedBlock, ok := envelope.Packet.P2PMessage.(*p2p.SignedBlock)
	if ok && signedBlock != nil {
		h.chain.PushBlock(chain.NewBlockStateByBlock(signedBlock))
	}
}

func (h *chainP2PHandler) Name() string {
	return "chain-p2p"
}
