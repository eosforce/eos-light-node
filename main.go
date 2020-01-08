package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eosforce/eos-light-node/core/chain"
	"github.com/eosforce/eos-light-node/eosio"
	"github.com/eosforce/eos-p2p/p2p"
	"github.com/eosforce/eos-p2p/store"
	"go.uber.org/zap"
)

var chainID = flag.String("chain-id", "bd61ae3a031e8ef2f97ee3b0e62776d6d30d4833c8f7c1645c657b149151004b", "net chainID to connect to")
var showLog = flag.Bool("v", false, "show detail log")
var startNum = flag.Int("num", 1, "start block num to sync")
var p2pAddress = flag.String("p2p", "", "p2p address")
var genesisPath = flag.String("genesis", "./config/genesis.json", "genesis file path")

var genesis *chain.Genesis

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
		//p2p.EnableP2PLogging()
	}

	var err error

	genesis, err = eosio.NewGenesisFromFile(*genesisPath)
	if err != nil {
		logger.Error("load genesis err", zap.Error(err))
		return
	}

	// from 9001 - 9020
	const maxNumListen int = 5
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

	peersCfg := make([]*p2p.PeerCfg, 0, len(peers))
	for _, p := range peers {
		peersCfg = append(peersCfg, &p2p.PeerCfg{
			Address: p,
		})
	}

	storer, err := store.NewBBoltStorer(logger, *chainID, "./blocks.db", false)
	if err != nil {
		logger.Error("new storer error", zap.Error(err))
		return
	}

	client, err := p2p.NewClient(
		ctx,
		*chainID,
		peersCfg,
		p2p.WithNeedSync(1),
		p2p.WithLogger(logger),
		p2p.WithStorer(storer),
		p2p.WithHandler(&chainP2PHandler{
			chain: chains,
		}),
	)

	if err != nil {
		logger.Error("p2p start error", zap.Error(err))
		return
	}

	// wait close node
	waitClose()

	// cancel context
	logger.Info("start close node")
	cancelFunc()

	logger.Info("wait p2p peers closed")
	client.Wait()

	logger.Info("wait chain closed")
	chains.Wait()

	logger.Info("wait storer closed")
	storer.Close()
	storer.Wait()

	logger.Info("light node closed success")
}
