package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fanyang1988/eos-light-node/p2p"
	"go.uber.org/zap"
)

var chainID = flag.String("chain-id", "1c6ae7719a2a3b4ecb19584a30ff510ba1b6ded86e1fd8b8fc22f1179c622a32", "net chainID to connect to")
var showLog = flag.Bool("v", true, "show detail log")
var startNum = flag.Int("num", 1, "start block num to sync")
var p2pAddress = flag.String("p2p", "", "p2p address")

// Wait wait for term signal, then stop the server
func Wait() {
	stopSignalChan := make(chan os.Signal, 1)
	signal.Notify(stopSignalChan,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		syscall.SIGUSR1)
	<-stopSignalChan
}

var (
	logger = zap.NewNop()
)

func main() {
	flag.Parse()

	if *showLog {
		logger, _ = zap.NewDevelopment()
	}

	// from 9001 - 9020
	const maxNumListen int = 1
	peers := make([]string, 0, maxNumListen+1)

	if *p2pAddress == "" {
		for i := 0; i < maxNumListen; i++ {
			peers = append(peers, fmt.Sprintf("127.0.0.1:%d", 9001+i))
		}
	} else {
		peers = append(peers, *p2pAddress)
	}

	logger.Sugar().Infof("start %v", *startNum)

	p2pPeers := p2p.NewP2PClient("p2p-peer", *chainID, 1, peers, logger)

	p2pPeers.RegHandler(p2p.NewHandlerLog(logger))
	p2pPeers.RegHandler(startHandler())
	err := p2pPeers.Start()

	if err != nil {
		logger.Error("start err", zap.Error(err))
	}

	Wait()

	err = p2pPeers.CloseConnection()
	if err != nil {
		logger.Error("start err", zap.Error(err))
	}
}

func startHandler() *p2p.HandlerToChannel {
	channel := make(chan p2p.MsgToChan, 1024)
	logger.Info("start handler")
	go func(ch chan p2p.MsgToChan) {
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
		}
	}(channel)

	return p2p.NewHandlerToChannel(channel)
}
