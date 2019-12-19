package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	chainp2p "github.com/fanyang1988/eos-p2p/p2p"
	"github.com/fanyang1988/eos-light-node/core/chain"
	"github.com/fanyang1988/eos-light-node/eosforce"
	"github.com/fanyang1988/eos-light-node/p2p"
	"go.uber.org/zap"
)

var chainID = flag.String("chain-id", "76eab2b704733e933d0e4eb6cc24d260d9fbbe5d93d760392e97398f4e301448", "net chainID to connect to")
var showLog = flag.Bool("v", false, "show detail log")
var startNum = flag.Int("num", 1, "start block num to sync")
var p2pAddress = flag.String("p2p", "", "p2p address")
var genesisPath = flag.String("genesis", "./config/genesis.json", "genesis file path")

var genesis *eosforce.Genesis

// waitClose wait for term signal, then stop the server
func waitClose() {
	stopSignalChan := make(chan os.Signal, 1)
	signal.Notify(stopSignalChan,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		syscall.SIGUSR1)
	<-stopSignalChan
}

func main() {
	flag.Parse()

	if *showLog {
		logger = newLogger(false)
		chainp2p.EnableP2PLogging()
	}

	var err error

	genesis, err = eosforce.NewGenesisFromFile(*genesisPath)
	if err != nil {
		logger.Error("load genesis err", zap.Error(err))
		return
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

	ctx, cancelFunc := context.WithCancel(context.Background())

	// TODO: now it is just for test
	chains := chain.New(ctx, logger)
	if err := chains.Init(genesis); err != nil {
		logger.Error("chains init error", zap.Error(err))
		return
	}

	p2pPeers := p2p.NewP2PClient(ctx, "p2p-peer", *chainID, 1, peers, logger)
	p2pPeers.RegHandler(p2p.NewHandlerLog(logger))
	p2pPeers.RegHandler(startHandler(ctx, chains))

	err = p2pPeers.Start(ctx)
	if err != nil {
		logger.Error("start err", zap.Error(err))
	}

	// wait close node
	waitClose()

	// cancel context
	logger.Info("start close node")
	cancelFunc()

	logger.Info("wait p2p peers closed")
	p2pPeers.Wait()

	logger.Info("wait chain closed")
	chains.Wait()

	logger.Info("light node closed success")
}
