package main

import (
	"github.com/fanyang1988/eos-light-node/p2p"
	"github.com/fanyang1988/eos-light-node/verifier"
	"go.uber.org/zap"
)

func startHandler() *p2p.HandlerToChannel {
	channel := make(chan p2p.MsgToChan, 1024)
	logger.Info("start handler")
	go func(ch chan p2p.MsgToChan) {
		verifier := verifier.NewBlockVerifier(genesis, logger)
		for {
			msg, ok := <-ch
			if !ok {
				logger.Error("handler chan close")
				return
			}

			if msg.CloseReason != 0 {
				logger.Error("handler chan close", zap.Uint8("reason", msg.CloseReason-1))
				return
			}

			//log.Logger().Info("handler block", zap.String("block", msg.Block.String()))
			err := verifier.Verify(&msg.Block)
			if err != nil {
				logger.Error("verify error", zap.Error(err))
			}
		}
	}(channel)

	return p2p.NewHandlerToChannel(channel)
}
