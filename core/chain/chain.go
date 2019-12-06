package chain

import "go.uber.org/zap"

// Chain eosc chain
type Chain struct {
	logger       *zap.Logger
	PendingState pendingState `json:"pending"`
}

func New(logger *zap.Logger) *Chain {
	return &Chain{
		logger: logger,
	}
}

// PushBlock try to append a block from net to chain,
// in eosio a block produced by self also need push block, but for a light node, cannot be producer
func (c *Chain) PushBlock(b *BlockState) error {
	// TODO use maybefork to select block status in chain
	return c.applyBlock(b, blockStatusIrreversible)
}

// startBlock init to ready to apply a block
func (c *Chain) startBlock(b *BlockState, blockState blockStatus) error {
	c.logger.Debug("startBlock",
		zap.Uint32("num", b.BlockNum), zap.String("id", b.BlockID))
	return nil
}

// finalizeBlock calc block id from chain and block base info, update chain status
func (c *Chain) finalizeBlock() error {
	return nil
}

// applyBlock apply a block with block state in chain to chain, call by PushBlock and Replay
func (c *Chain) applyBlock(b *BlockState, blockState blockStatus) error {
	if err := c.startBlock(b, blockState); err != nil {
		return err
	}

	if err := c.finalizeBlock(); err != nil {
		return err
	}

	return nil
}
