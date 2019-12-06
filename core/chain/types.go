package chain

import (
	eos "github.com/eosforce/goeosforce"
	"github.com/eosforce/goeosforce/ecc"
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
type SHA256Bytes = Checksum256

// IncrementalMerkle for block root Merkle
type IncrementalMerkle = types.IncrementalMerkle

// ProducerKey eos ProducerKey type
type ProducerKey = eos.ProducerKey
type PublicKey = ecc.PublicKey
type PrivateKey = ecc.PrivateKey

// ProducerSchedule eos ProducerSchedule type
type ProducerSchedule = eos.ProducerSchedule

type AccountName = eos.AccountName
type PermissionName = eos.PermissionName
type ActionName = eos.ActionName
type TableName = eos.TableName
type ScopeName = eos.ScopeName

func AN(in string) AccountName    { return AccountName(in) }
func ActN(in string) ActionName   { return ActionName(in) }
func PN(in string) PermissionName { return PermissionName(in) }

func MarshalBinary(v interface{}) ([]byte, error) {
	return eos.MarshalBinary(v)
}

func MustNewPublicKey(pubKey string) PublicKey {
	return ecc.MustNewPublicKey(pubKey)
}

// for p2p
type Packet = eos.Packet

const GoAwayMessageType = eos.GoAwayMessageType

type GoAwayMessage = eos.GoAwayMessage

const SignedBlockType = eos.SignedBlockType
