package verifier

import (
	"github.com/eoscanada/eos-go"
	"github.com/eosspark/eos-go/chain/types"
	"go.uber.org/zap"
)

type BlockHeaderStatus struct {
	logger           *zap.Logger
	BlockID          eos.Checksum256         `json:"block_id"`
	BlockNum         uint32                  `json:"block_num"`
	Previous         eos.Checksum256         `json:"previous"`
	PreviousBlockNum uint32                  `json:"previous_num"`
	BlockrootMerkle  types.IncrementalMerkle `json:"blockroot_merkle"`
}

func (b *BlockHeaderStatus) ToNext(block *eos.SignedBlock) {
	id, _ := block.BlockID()
	b.BlockrootMerkle.Append(ToSha256(id))
	b.Previous = id
	b.BlockNum = block.BlockNumber() + 1

	b.logger.Debug("to next",
		zap.Uint32("num", b.BlockNum),
		zap.String("previous", b.Previous.String()),
		zap.String("block root merkle", b.BlockrootMerkle.GetRoot().String()))
}
