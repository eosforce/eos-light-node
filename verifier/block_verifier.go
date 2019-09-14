package verifier

import (
	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

// BlockVerifier block verifier
type BlockVerifier struct {
	logger    *zap.Logger
	producers ScheduleProducersDatas
}

// NewBlockVerifier create a new BlockVerifier
func NewBlockVerifier(logger *zap.Logger) *BlockVerifier {
	return &BlockVerifier{
		logger: logger,
		producers: ScheduleProducersDatas{
			logger:    logger,
			schedules: make([]scheduleProducers, 4096),
		},
	}
}

// Verify verifier block
func (v *BlockVerifier) Verify(block *eos.SignedBlock) error {
	v.logger.Debug("verify block", zap.Uint32("number", block.BlockNumber()))

	err := v.producers.OnBlock(block)
	if err != nil {
		return err
	}

	return nil
}
