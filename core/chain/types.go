package chain

import (
	eos "github.com/eosforce/goeosforce"
	"github.com/eosforce/goeosforce/ecc"
	chaintype "github.com/fanyang1988/eos-light-node/eosforce"
)

// SignedBlock Signed block in chain, for a light node, all block will be signed from others
type SignedBlock = eos.SignedBlock

// BlockState block detail state from a signed block data and the chain state
type BlockState = eos.BlockState

// BlockHeader header data for block
type BlockHeader = eos.BlockHeader

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

// SHA256Bytes for sha256 from eosio lib
type SHA256Bytes = Checksum256

// ProducerKey eos ProducerKey type
type ProducerKey = eos.ProducerKey

// PublicKey ecc.PublicKey
type PublicKey = ecc.PublicKey

// PrivateKey ecc.PrivateKey
type PrivateKey = ecc.PrivateKey

// ProducerSchedule eos ProducerSchedule type
type ProducerSchedule = eos.ProducerSchedule

// AccountName eos.AccountName
type AccountName = eos.AccountName

// PermissionName eos.PermissionName
type PermissionName = eos.PermissionName

// ActionName eos.ActionName
type ActionName = eos.ActionName

// TableName eos.TableName
type TableName = eos.TableName

// ScopeName eos.ScopeName
type ScopeName = eos.ScopeName

// AN from string to account name
func AN(in string) AccountName { return AccountName(in) }

// ActN from string to action name
func ActN(in string) ActionName { return ActionName(in) }

// PN from string to permission name
func PN(in string) PermissionName { return PermissionName(in) }

// MarshalBinary call eos MarshalBinary
func MarshalBinary(v interface{}) ([]byte, error) {
	return eos.MarshalBinary(v)
}

// MustNewPublicKey call ecc MustNewPublicKey
func MustNewPublicKey(pubKey string) PublicKey {
	return ecc.MustNewPublicKey(pubKey)
}

// TypeSize size for eos types
var TypeSize = eos.TypeSize

// for p2p

// Packet eos.Packet
type Packet = eos.Packet

// GoAwayMessage eos.GoAwayMessage
type GoAwayMessage = eos.GoAwayMessage

const (
	// SignedBlockType eos.SignedBlockType
	SignedBlockType = eos.SignedBlockType
	// GoAwayMessageType eos.GoAwayMessageType
	GoAwayMessageType = eos.GoAwayMessageType
)

// Genesis genesis datas for chain
type Genesis = chaintype.Genesis

// Action chain action
type Action = eos.Action

// Transaction chain transaction
type Transaction = eos.Transaction
