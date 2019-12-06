package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fanyang1988/eos-light-node/eosforce"
	"github.com/fanyang1988/eos-light-node/p2p"
	"go.uber.org/zap"
)

var chainID = flag.String("chain-id", "a1c3bfe884d9cad5dcd496a09bca555771fd4e6b1fea164542018482b39ea3f4", "net chainID to connect to")
var showLog = flag.Bool("v", false, "show detail log")
var startNum = flag.Int("num", 1, "start block num to sync")
var p2pAddress = flag.String("p2p", "", "p2p address")
var genesisPath = flag.String("genesis", "./config/genesis.json", "genesis file path")

var genesis *eosforce.Genesis

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

func main() {
	flag.Parse()

	if *showLog {
		logger = newLogger(false)
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

	p2pPeers := p2p.NewP2PClient("p2p-peer", *chainID, 1, peers, logger)
	p2pPeers.RegHandler(p2p.NewHandlerLog(logger))
	p2pPeers.RegHandler(startHandler())
	err = p2pPeers.Start()

	if err != nil {
		logger.Error("start err", zap.Error(err))
	}

	Wait()

	err = p2pPeers.CloseConnection()
	if err != nil {
		logger.Error("start err", zap.Error(err))
	}
}
