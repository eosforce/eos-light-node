package chain

// Chain eosc chain
type Chain struct {
}

// PushBlock try to append a block from net to chain,
// in eosio a block produced by self also need push block, but for a light node, cannot be producer
func (c *Chain) PushBlock(b *BlockState) error {
	return nil
}

// startBlock init to ready to apply a block
func (c *Chain) startBlock(b *BlockState, blockState blockStatus) error {
	return nil
}

// finalizeBlock calc block id from chain and block base info, update chain status
func (c *Chain) finalizeBlock() error {
	return nil
}

// applyBlock apply a block with block state in chain to chain, call by PushBlock and Replay
func (c *Chain) applyBlock(b *BlockState, blockState blockStatus) error {
	return nil
}
