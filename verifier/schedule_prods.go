package verifier

import (
	"crypto/sha256"
	"errors"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"go.uber.org/zap"
)

type scheduleProducers struct {
	version   uint32
	blockNum  uint32
	hash      eos.Checksum256
	producers []eos.ProducerKey
}

// ScheduleProducersDatas all schedule producers datas by version
type ScheduleProducersDatas struct {
	logger    *zap.Logger
	schedules []scheduleProducers
}

func (s *ScheduleProducersDatas) appendDatas(sp *scheduleProducers) error {
	if sp.version != uint32(len(s.schedules)) {
		s.logger.Error("too early schedule version", zap.Uint32("version", sp.version), zap.Int("len", len(s.schedules)))
		return errors.New("too early schedule version")
	}

	s.schedules = append(s.schedules, *sp)

	s.logger.Info("new schedule producers version",
		zap.Uint32("version", sp.version), zap.String("hash", sp.hash.String()))

	return nil
}

func (s *ScheduleProducersDatas) Init() {
	s.OnNewProducers(1, eos.ProducerSchedule{
		Version: 0,
		Producers: []eos.ProducerKey{
			eos.ProducerKey{
				AccountName:     eos.AN("eosio"),
				BlockSigningKey: ecc.MustNewPublicKey("EOS8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr"), // TODO: use genesis data
			},
		},
	})
}

func (s *ScheduleProducersDatas) OnNewProducers(blockNum uint32, n eos.ProducerSchedule) error {
	spRaws, err := eos.MarshalBinary(n)
	if err != nil {
		return err
	}

	h := sha256.New()
	_, _ = h.Write(spRaws)

	sp := &scheduleProducers{
		version:   n.Version,
		blockNum:  blockNum,
		producers: make([]eos.ProducerKey, 0, 30),
		hash:      h.Sum(nil),
	}

	for _, p := range n.Producers {
		sp.producers = append(sp.producers, p)
	}

	return s.appendDatas(sp)
}

// GetScheduleProducer get producer by version and account name
func (s *ScheduleProducersDatas) GetScheduleProducer(version uint32, name eos.AccountName) (eos.ProducerKey, error) {
	if version >= uint32(len(s.schedules)) {
		return eos.ProducerKey{}, errors.New("no version found")
	}

	for _, sp := range s.schedules[version].producers {
		if sp.AccountName == name {
			return sp, nil
		}
	}

	return eos.ProducerKey{}, errors.New("no producer in version datas")
}

func (s *ScheduleProducersDatas) GetScheduleProducersHash() eos.Checksum256 {
	return s.schedules[len(s.schedules)-1].hash // must has value
}

// OnBlock on block to update datas
func (s *ScheduleProducersDatas) OnBlock(msg *eos.SignedBlock) error {
	if msg.NewProducers == nil {
		return nil
	}

	return s.OnNewProducers(msg.BlockNumber(), msg.NewProducers.ProducerSchedule)
}
