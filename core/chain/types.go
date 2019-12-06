package chain

import (
	eos "github.com/eosforce/goeosforce"
	"github.com/eosspark/eos-go/chain/types"
)

// SignedBlock Signed block in chain, for a light node, all block will be signed from others
type SignedBlock = eos.SignedBlock

// BlockState block detail state from a signed block data and the chain state
type BlockState = eos.BlockState

// NewBlockStateByBlock just tmp imp
func NewBlockStateByBlock(sb *SignedBlock) *BlockState {
	b := &BlockState{}

	blockID, _ := sb.BlockID()

	b.BlockID = blockID.String()
	b.BlockNum = sb.BlockNumber()

	// TODO:
	b.DPoSProposedIrreversibleBlockNum = b.BlockNum
	b.DPoSIrreversibleBlockNum = b.BlockNum
	b.ActiveSchedule = nil
	b.BlockrootMerkle = nil

	b.SignedBlock = sb

	return b
}

// Checksum256 id
type Checksum256 = eos.Checksum256

// IncrementalMerkle for block root Merkle
type IncrementalMerkle = types.IncrementalMerkle
