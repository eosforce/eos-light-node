package chain

import "errors"

var (
	// ErrChainFork error by block fork in chain
	ErrChainFork = errors.New("errorChainFork")
)
