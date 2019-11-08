package chain

import (
	eos "github.com/eosforce/goeosforce"
)

// SignedBlock Signed block in chain, for a light node, all block will be signed from others
type SignedBlock = eos.SignedBlock

// BlockState block detail state from a signed block data and the chain state
type BlockState = eos.BlockState
