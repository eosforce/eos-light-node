package p2p

import (
	eos "github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

// p2pHandler handler for p2p client
type p2pHandler interface {
	OnBlock(peer string, msg *eos.SignedBlock) error
	OnGoAway(peer string, reason uint8, nodeID eos.Checksum256) error
}

// HandlerToLog p2p msg handler to log msg received
type HandlerToLog struct {
	logger *zap.Logger
}

// NewHandlerLog create HandlerToLog handler
func NewHandlerLog(l *zap.Logger) *HandlerToLog {
	return &HandlerToLog{
		logger: l,
	}
}

// OnBlock log block received
func (h HandlerToLog) OnBlock(peer string, msg *eos.SignedBlock) error {
	blockID, err := msg.BlockID()
	if err != nil {
		return err
	}
	h.logger.Sugar().Infof("OnBlock %s - [%d]: id: %s, sv: %d, trx: %d.",
		peer, msg.BlockNumber(), blockID,
		msg.ScheduleVersion,
		len(msg.Transactions))

	if msg.NewProducers != nil {
		h.logger.Sugar().Infof("NewProducers %s - [%d]: %d", peer, msg.BlockNumber(), msg.NewProducers.Version)
		for idx, producer := range msg.NewProducers.Producers {
			h.logger.Sugar().Infof("NewProducers %s - [%d]: producer %d: %s - %s",
				peer, msg.BlockNumber(),
				idx, producer.AccountName, producer.BlockSigningKey)
		}
	}
	return nil
}

// OnGoAway log goaway message received
func (h HandlerToLog) OnGoAway(peer string, reason uint8, nodeID eos.Checksum256) error {
	h.logger.Sugar().Errorf("OnGoAway %s by %d", peer, reason)
	return nil
}
