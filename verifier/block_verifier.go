package verifier

import (
	"errors"

	eos "github.com/eosforce/goeosforce"
	"github.com/fanyang1988/eos-light-node/eosforce"
	"go.uber.org/zap"
)

// BlockVerifier block verifier
type BlockVerifier struct {
	logger    *zap.Logger
	producers ScheduleProducersDatas
	status    BlockHeaderStatus
	genesis   eosforce.Genesis
}

// NewBlockVerifier create a new BlockVerifier
func NewBlockVerifier(genesis *eosforce.Genesis, logger *zap.Logger) *BlockVerifier {
	res := &BlockVerifier{
		logger: logger,
		producers: ScheduleProducersDatas{
			logger:    logger,
			schedules: make([]scheduleProducers, 0, 4096),
		},
		status: BlockHeaderStatus{
			logger:   logger,
			BlockNum: 1,
		},
		genesis: *genesis,
	}
	res.producers.Init(&res.genesis)
	return res
}

func (v *BlockVerifier) GetSigDigest(block *eos.SignedBlock) (eos.Checksum256, error) {
	headerHash := GetBlockHeaderHash(&block.BlockHeader)
	scheduleProducersHash := v.producers.GetScheduleProducersHash()
	blockrootMerkle := v.status.BlockrootMerkle.GetRoot()

	headerAndBmroot := HashCheckSumPairH256(headerHash, blockrootMerkle)
	sigDigest := HashCheckSumPair(headerAndBmroot, scheduleProducersHash)

	/*
		v.logger.Debug("sigDigest info",
			zap.String("headerHash", headerHash.String()),
			zap.String("scheduleProducersHash", scheduleProducersHash.String()),
			zap.String("headerAndBmroot", headerAndBmroot.String()),
			zap.String("sigDigest", sigDigest.String()))
	*/
	return sigDigest, nil
}

func (v *BlockVerifier) VerifySign(block *eos.SignedBlock) error {
	sigDigest, err := v.GetSigDigest(block)
	if err != nil {
		return err
	}

	pubKey, err := v.producers.GetScheduleProducer(block.ScheduleVersion, block.Producer)
	if err != nil {
		return err
	}

	signPubKey, err := block.ProducerSignature.PublicKey(sigDigest)
	if err != nil {
		return err
	}

	if !IsSamePubKey(pubKey.BlockSigningKey, signPubKey) {
		v.logger.Error("sign err",
			zap.String("pubkey", pubKey.BlockSigningKey.String()),
			zap.String("signPubKey", signPubKey.String()),
			zap.Uint32("scheduleVersion", block.ScheduleVersion),
			zap.String("producer", string(block.Producer)),
		)
		return errors.New("sign error")
	}

	return nil
}

// Verify verifier block
func (v *BlockVerifier) Verify(block *eos.SignedBlock) error {
	blockNum := block.BlockNumber()

	if blockNum != v.status.BlockNum {
		return nil
	}

	if blockNum != 1 {
		err := v.producers.OnBlock(block)
		if err != nil {
			return err
		}

		err = v.VerifySign(block)
		if err != nil {
			return err
		}

		if blockNum%1000 == 0 {
			v.logger.Info("verify block", zap.Uint32("number", blockNum))
		}
	}

	v.status.ToNext(block)
	return nil
}
