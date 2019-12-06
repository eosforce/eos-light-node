package p2p

import (
	"encoding/hex"

	"github.com/pkg/errors"

	"github.com/fanyang1988/eos-light-node/core/chain"
)

// Config config to codex-go
type Config struct {
	ChainID chain.SHA256Bytes
	URL     string
	Keys    map[string]accountKey
	Prikeys []chain.PrivateKey
	IsDebug bool
}

type accountKey struct {
	Name   chain.AccountName
	PubKey chain.PublicKey
	PriKey chain.PrivateKey
}

// ToSHA256Bytes from string to sha256
func ToSHA256Bytes(in string) (chain.SHA256Bytes, error) {
	if len(in) != 64 {
		return nil, errors.New("should be 64 hexadecimal characters")
	}

	bytes, err := hex.DecodeString(in)
	if err != nil {
		return nil, err
	}

	return chain.SHA256Bytes(bytes), nil
}
