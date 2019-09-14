package verifier

import (
	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

// BlockVerifier block verifier
type BlockVerifier struct {
	logger *zap.Logger
}

// NewBlockVerifier create a new BlockVerifier
func NewBlockVerifier(logger *zap.Logger) *BlockVerifier {
	return &BlockVerifier{
		logger: logger,
	}
}

// Verify verifier block
func (v *BlockVerifier) Verify(block *eos.SignedBlock) error {
	v.logger.Debug("verify block", zap.Uint32("number", block.BlockNumber()))
	return nil
}
