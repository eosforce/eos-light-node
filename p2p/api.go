package p2p

import "go.uber.org/zap"

type P2PInitParams struct {
	Name          string   `json:"name"`
	ClientID      string   `json:"clientID"`
	Peers         []string `json:"peers"`
	StartBlockNum uint32   `json:"start"`
	Logger        *zap.Logger
}
