package p2p

import (
	"encoding/hex"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/pkg/errors"
)

// Config config to codex-go
type Config struct {
	ChainID eos.SHA256Bytes
	URL     string
	Keys    map[string]accountKey
	Prikeys []ecc.PrivateKey
	IsDebug bool
}

type accountKey struct {
	Name   eos.AccountName
	PubKey ecc.PublicKey
	PriKey ecc.PrivateKey
}

// ToSHA256Bytes from string to sha256
func ToSHA256Bytes(in string) (eos.SHA256Bytes, error) {
	if len(in) != 64 {
		return nil, errors.New("should be 64 hexadecimal characters")
	}

	bytes, err := hex.DecodeString(in)
	if err != nil {
		return nil, err
	}

	return eos.SHA256Bytes(bytes), nil
}
