package chain

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Chain eosc chain
type Chain struct {
	logger    *zap.Logger
	publisher *messagePublisher

	PendingState      pendingState           `json:"pending"`
	ScheduleProducers ScheduleProducersDatas `json:"scheduleProducers"`
}

// New create chain
func New(ctx context.Context, logger *zap.Logger) *Chain {
	return &Chain{
		logger:    logger,
		publisher: newMessagePublisher(ctx, logger),
		ScheduleProducers: ScheduleProducersDatas{
			logger:    logger,
			schedules: make([]scheduleProducers, 0, 4096),
		},
	}
}

// Init init the chain
func (c *Chain) Init(genesis *Genesis) error {
	c.ScheduleProducers.Init(genesis)
	c.PendingState.BlockNum = 1
	return nil
}

func (c *Chain) Wait() {
	c.publisher.Wait()
}

// PushBlock try to append a block from net to chain,
// in eosio a block produced by self also need push block, but for a light node, cannot be producer
func (c *Chain) PushBlock(b *BlockState) error {
	// TODO use maybefork to select block status in chain
	if c.PendingState.BlockNum != b.BlockNum {
		return nil
	}
	return c.applyBlock(b, blockStatusIrreversible)
}

// startBlock init to ready to apply a block
func (c *Chain) startBlock(b *BlockState, blockState blockStatus) error {
	//c.logger.Debug("startBlock",
	//	zap.Uint32("num", b.BlockNum), zap.String("id", b.BlockID))

	// update pending data, init building block data

	// check preactivated_features (  no need in current light node  )

	// update schedule datas
	if err := c.ScheduleProducers.OnBlock(b.SignedBlock); err != nil {
		return errors.Wrapf(err, "update schedule datas err in %d", b.BlockNum)
	}

	// onblock ( no need in current light node )

	// clear_expired_input_transactions ( no need in current light node )

	// update_producers_authority --> update eosio.prods auth ( no need in current light node )

	return nil
}

// finalizeBlock calc block id from chain and block base info, update chain status
func (c *Chain) finalizeBlock(b *BlockState) error {

	// verify block is signed
	if b.BlockNum != 1 {
		if err := c.verifySign(b.SignedBlock); err != nil {
			return errors.Wrapf(err, "verify sign error %d", b.BlockNum)
		}
	}

	// check if block is building ok in pending

	// update resource limits

	// create a unsigned block with building block data in pending( which is not getted block )
	// this data is to build assembled_block data

	// Update TaPoS table by create_block_summary

	// update pending with assembled_block data, which is gen by building block data
	c.PendingState.update(b.SignedBlock)

	c.publisher.OnCommittedBlock(b.SignedBlock)

	return nil
}

// applyBlock apply a block with block state in chain to chain, call by PushBlock and Replay
func (c *Chain) applyBlock(b *BlockState, blockState blockStatus) error {
	// start_block
	if err := c.startBlock(b, blockState); err != nil {
		return err
	}

	// push all trx, check trx exec res is same

	// finalize_block
	if err := c.finalizeBlock(b); err != nil {
		return err
	}

	// get exec res(assembled_block) in finalize_block, check id is same with block getted

	// update completed block info in pending

	return nil
}
