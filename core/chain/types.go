package chain

import (
	eos "github.com/eosforce/goeosforce"
	"github.com/eosspark/eos-go/chain/types"
)

// SignedBlock Signed block in chain, for a light node, all block will be signed from others
type SignedBlock = eos.SignedBlock

// BlockState block detail state from a signed block data and the chain state
type BlockState = eos.BlockState

// Checksum256 id
type Checksum256 = eos.Checksum256

// IncrementalMerkle for block root Merkle
type IncrementalMerkle = types.IncrementalMerkle
