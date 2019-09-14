package verifier

import (
	"errors"

	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

type scheduleProducers struct {
	version   uint32
	blockNum  uint32
	producers []eos.ProducerKey
}

// ScheduleProducersDatas all schedule producers datas by version
type ScheduleProducersDatas struct {
	logger    *zap.Logger
	schedules []scheduleProducers
}

func (s *ScheduleProducersDatas) appendDatas(sp *scheduleProducers) error {
	if sp.version != uint32(len(s.schedules)+1) {
		return errors.New("too early schedule version")
	}

	s.schedules = append(s.schedules, *sp)

	s.logger.Info("new schedule producers version", zap.Uint32("version", sp.version))

	return nil
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

// OnBlock on block to update datas
func (s *ScheduleProducersDatas) OnBlock(msg *eos.SignedBlock) error {
	if msg.NewProducers == nil {
		return nil
	}

	sp := &scheduleProducers{
		version:   msg.NewProducers.Version,
		blockNum:  msg.BlockNumber(),
		producers: make([]eos.ProducerKey, 0, 30),
	}

	for _, p := range msg.NewProducers.Producers {
		sp.producers = append(sp.producers, p)
	}

	return s.appendDatas(sp)

}
