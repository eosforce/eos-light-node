package chain

// pendingState pending start in apply chain, include last block root hash
type pendingState struct {
	BlockStatus      blockStatus       `json:"block_status"`
	BlockID          Checksum256       `json:"block_id"`
	BlockNum         uint32            `json:"block_num"`
	Previous         Checksum256       `json:"previous"`
	PreviousBlockNum uint32            `json:"previous_num"`
	BlockrootMerkle  IncrementalMerkle `json:"blockroot_merkle"`
}

func (p *pendingState) update(block *SignedBlock) {
	id, _ := block.BlockID()
	p.BlockrootMerkle.Append(id)
	p.Previous = id
	p.BlockNum = block.BlockNumber() + 1
}
