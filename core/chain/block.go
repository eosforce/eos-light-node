package chain

// blockStatus block status in chain
type blockStatus uint8

const (
	// blockStatusIrreversible this block has already been applied before by this node and is considered irreversible
	blockStatusIrreversible blockStatus = iota
	// blockStatusValidated this is a complete block signed by a valid producer and has been previously applied by this node and therefore validated but it is not yet irreversible
	blockStatusValidated
	// blockStatusComplete this is a complete block signed by a valid producer but is not yet irreversible nor has it yet been applied by this node
	blockStatusComplete
	// blockStatusIncomplete this is an incomplete block (either being produced by a producer or speculatively produced by a node)
	blockStatusIncomplete
)
