package chain

import "time"

// ActionMsg msg when a action is committed to the chain
type ActionMsg struct {
	BlockNum      uint32
	BlockID       Checksum256
	TransactionID Checksum256
	Timestamp     time.Time
	Act           Action
}

type actMsgChan chan ActionMsg

// TransactionMsg msg when a transaction is committed to the chain
type TransactionMsg struct {
	BlockNum  uint32
	BlockID   Checksum256
	Timestamp time.Time
	Trx       Transaction
}

type trxMsgChan chan TransactionMsg

// BlockMsg msg when a block is committed to the chain
type BlockMsg struct {
	Block SignedBlock
}

type blockMsgChan chan BlockMsg

// ActMsgHandler act msg handler func type
type ActMsgHandler func(act *ActionMsg)

// TrxMsgHandler trx msg handler func type
type TrxMsgHandler func(act *TransactionMsg)

// BlkMsgHandler blk msg handler func type
type BlkMsgHandler func(act *BlockMsg)

func (c *Chain) startHandler() {

}

func (c *Chain) onBlockCommit(b *SignedBlock) {

}
