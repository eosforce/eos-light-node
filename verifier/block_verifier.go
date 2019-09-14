package verifier

import (
	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

// BlockVerifier block verifier
type BlockVerifier struct {
	logger    *zap.Logger
	producers ScheduleProducersDatas
	status    BlockHeaderStatus
}

// NewBlockVerifier create a new BlockVerifier
func NewBlockVerifier(logger *zap.Logger) *BlockVerifier {
	return &BlockVerifier{
		logger: logger,
		producers: ScheduleProducersDatas{
			logger:    logger,
			schedules: make([]scheduleProducers, 0, 4096),
		},
		status: BlockHeaderStatus{
			logger:   logger,
			BlockNum: 1,
		},
	}
}

// Verify verifier block
func (v *BlockVerifier) Verify(block *eos.SignedBlock) error {
	blockNum := block.BlockNumber()

	if blockNum != v.status.BlockNum {
		return nil
	}

	//if blockNum%1000 == 0 {
	v.logger.Info("verify block", zap.Uint32("number", blockNum))
	//}

	err := v.producers.OnBlock(block)
	if err != nil {
		return err
	}

	v.status.ToNext(block)
	return nil
}
