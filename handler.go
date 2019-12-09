package main

import (
	"github.com/fanyang1988/eos-light-node/core/chain"
	"github.com/fanyang1988/eos-light-node/eosforce"
	"github.com/fanyang1988/eos-light-node/p2p"
	"go.uber.org/zap"
)

func startHandler() *p2p.HandlerToChannel {
	channel := make(chan p2p.MsgToChan, 1024)
	logger.Info("start handler")
	go func(ch chan p2p.MsgToChan) {

		logger.Sugar().Infof("init handler goroutine")

		// TODO: now it is just for test
		chains := chain.New(logger)
		genesis, _ := eosforce.NewGenesisFromFile("./genesis.json")
		chains.Init(genesis)

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

			logger.Info("handler block", zap.String("block", msg.Block.String()))
			err := chains.PushBlock(chain.NewBlockStateByBlock(&msg.Block))
			if err != nil {
				logger.Error("push block error", zap.Error(err))
				panic(err)
			}
		}
	}(channel)

	return p2p.NewHandlerToChannel(channel)
}
