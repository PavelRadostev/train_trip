package listener

import (
	"github.com/as-master/train_trip/internal/domain"
	cqrs "github.com/as-master/train_trip/internal/domain/cqrs"
)

// TrainLoadQuery — запрос на загрузку поезда
type TrainLoadQuery struct {
	Request Request

	TrainID      []int           `cbor:"train_id,omitempty"`
	TimeInterval domain.Interval `cbor:"time_interval,omitempty"`
	ShovelID     int             `cbor:"shovel_id,omitempty"`
	UnloadID     int             `cbor:"unload_id,omitempty"`
	GeometryID   int             `cbor:"geometry_id,omitempty"`
	OnlyFinished bool            `cbor:"only_finished,omitempty"`
}

func (TrainLoadQuery) StreamKey() string {
	return "train_load_unload_detector_domain.query.train_load_unload.TrainLoadQuery"
}

func (TrainLoadQuery) Handle(msg string) string {
	return "train_load_unload_detector_domain.query.train_load_unload.TrainLoadQuery"
}

// TrainUnloadQuery — запрос на выгрузку поезда
type TrainUnloadQuery struct {
	Request Request

	TrainID      []int           `cbor:"train_id,omitempty"`
	TimeInterval domain.Interval `cbor:"time_interval,omitempty"`
	UnloadID     int             `cbor:"unload_id,omitempty"`
	GeometryID   int             `cbor:"geometry_id,omitempty"`
	OnlyFinished bool            `cbor:"only_finished,omitempty"`
}

func (TrainUnloadQuery) StreamKey() string {
	return "train_load_unload_detector_domain.query.train_load_unload.TrainUnloadQuery"
}

func (TrainUnloadQuery) Handle(msg string) string {
	return "train_load_unload_detector_domain.query.train_load_unload.TrainUnloadQuery"
}

// Зарегистрируем запросы в регистраторе создавая слушателей потоков
func RegisterQueries() *cqrs.Registrar {
	reg := cqrs.GetRegistrar()
	reg.Register(TrainLoadQuery{})
	reg.Register(TrainUnloadQuery{})

	return reg
}
