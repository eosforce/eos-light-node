package main

import (
	"context"

	"github.com/fanyang1988/eos-light-node/core/chain"
	"github.com/fanyang1988/eos-light-node/p2p"
	"go.uber.org/zap"
)

func startHandler(ctx context.Context, chains *chain.Chain) *p2p.HandlerToChannel {
	channel := make(chan p2p.MsgToChan, 1024)
	logger.Info("start handler")
	go func(ch chan p2p.MsgToChan) {
		logger.Sugar().Infof("init handler goroutine")
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					logger.Error("handler chan close")
					return
				}

				// TODO: stop by context
				if msg.CloseReason != 0 {
					logger.Error("handler chan close", zap.Uint8("reason", msg.CloseReason-1))
					return
				}

				logger.Info("handler block", zap.String("block", msg.Block.String()))
				err := chains.PushBlock(chain.NewBlockStateByBlock(&msg.Block))
				if err != nil {
					logger.Error("push block error", zap.Error(err))
					panic(err)
				}
			case <-ctx.Done():
				logger.Info("chains p2p handler close")
				return
			default:
			}
		}
	}(channel)

	return p2p.NewHandlerToChannel(channel)
}
