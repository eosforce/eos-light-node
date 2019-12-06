package p2p

import (
	"github.com/fanyang1988/eos-light-node/core/chain"
	"go.uber.org/zap"
)

// p2pHandler handler for p2p client
type p2pHandler interface {
	OnBlock(peer string, msg *chain.SignedBlock) error
	OnGoAway(peer string, reason uint8, nodeID chain.Checksum256) error
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
func (h HandlerToLog) OnBlock(peer string, msg *chain.SignedBlock) error {
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
func (h HandlerToLog) OnGoAway(peer string, reason uint8, nodeID chain.Checksum256) error {
	h.logger.Sugar().Errorf("OnGoAway %s by %d", peer, reason)
	return nil
}

// MsgToChan type chan
type MsgToChan struct {
	Peer        string
	CloseReason uint8 // if not zero, mean go away reason
	Block       chain.SignedBlock
}

// HandlerToChannel handler to send msg to a channel
type HandlerToChannel struct {
	channel   chan<- MsgToChan
	hasClosed bool
}

// NewHandlerToChannel create HandlerToChannel
func NewHandlerToChannel(channel chan<- MsgToChan) *HandlerToChannel {
	return &HandlerToChannel{
		channel:   channel,
		hasClosed: false,
	}
}

// OnBlock block received
func (h HandlerToChannel) OnBlock(peer string, msg *chain.SignedBlock) error {
	if h.hasClosed {
		return nil
	}

	h.channel <- MsgToChan{
		Peer:  peer,
		Block: *msg,
	}
	return nil
}

// OnGoAway goaway message received
func (h HandlerToChannel) OnGoAway(peer string, reason uint8, nodeID chain.Checksum256) error {
	if h.hasClosed {
		return nil
	}

	h.channel <- MsgToChan{
		Peer:        peer,
		CloseReason: reason + 1,
	}
	close(h.channel)
	h.hasClosed = true
	return nil
}
