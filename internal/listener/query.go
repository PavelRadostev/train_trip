package listener

import (
	"context"
	"fmt"

	"github.com/as-master/train_trip/internal/domain"
	"github.com/as-master/train_trip/internal/domain/model"
	"github.com/as-master/train_trip/pkg/cqrs"
)

// TrainLoadQuery — запрос на загрузку поезда
type TrainLoadQuery struct {
	TrainID      []int           `cbor:"train_id,omitempty"`
	TimeInterval domain.Interval `cbor:"time_interval,omitempty"`
	ShovelID     int             `cbor:"shovel_id,omitempty"`
	UnloadID     int             `cbor:"unload_id,omitempty"`
	GeometryID   int             `cbor:"geometry_id,omitempty"`
	OnlyFinished bool            `cbor:"only_finished,omitempty"`
}

func (q TrainLoadQuery) StreamKey() string {
	return "train_load_unload_detector_domain.query.train_load_unload.TrainLoadQuery"
}

func (q TrainLoadQuery) Handle(repo cqrs.Repository) (any, error) {
	ctx := context.TODO()
	conn := repo.GetConnection()
	query := `
		SELECT 
		id, is_deleted, load_arrive_time, load_begin_time, load_end_time,
		load_depart_time, train_id, geometry_id, unload_id, shovel_id,
		manual, load_type_id_manual, work_type_id_manual,
		load_begin_time_manual, load_end_time_manual,
		unload_id_manual, shovel_id_manual, volume_manual,
		cycle_ids, is_cured, carriage_num, source
		FROM train_load_unload_store.loads 
		WHERE train_id = ANY($1) 
		AND load_begin_time >= $2::timestamp 
		AND load_end_time <= $3::timestamp
	`

	rows, err := conn.Query(ctx, query, q.TrainID, q.TimeInterval.TimeFrom, q.TimeInterval.TimeTo)
	if err != nil {
		return nil, fmt.Errorf("train_load_unload_detector_domain.query.train_load_unload.TrainLoadQuery: failed to execute query: %w", err)
	}
	defer rows.Close()

	result := make([]model.TrainLoad, 0, 40) // 40 половина среднего количества погруок за смену
	for rows.Next() {
		var load model.TrainLoad
		if err := rows.Scan(
			&load.ID,
			&load.IsDeleted,
			&load.LoadArriveTime,
			&load.LoadBeginTime,
			&load.LoadEndTime,
			&load.LoadDepartTime,
			&load.TrainID,
			&load.GeometryID,
			&load.UnloadID,
			&load.ShovelID,
			&load.Manual,
			&load.LoadTypeIDManual,
			&load.WorkTypeIDManual,
			&load.LoadBeginTimeManual,
			&load.LoadEndTimeManual,
			&load.UnloadIDManual,
			&load.ShovelIDManual,
			&load.VolumeManual,
			&load.CycleIDs,
			&load.IsCured,
			&load.CarriageNum,
			&load.Source,
		); err != nil {
			return nil, fmt.Errorf("train_load_unload_detector_domain.query.train_load_unload.TrainLoadQuery: failed to scan row: %w", err)
		}
		result = append(result, load)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("train_load_unload_detector_domain.query.train_load_unload.TrainLoadQuery: row iteration error: %w", err)
	}
	return result, nil
}

// TrainUnloadQuery — запрос на выгрузку поезда
type TrainUnloadQuery struct {
	TrainID      []int           `cbor:"train_id,omitempty"`
	TimeInterval domain.Interval `cbor:"time_interval,omitempty"`
	UnloadID     int             `cbor:"unload_id,omitempty"`
	GeometryID   int             `cbor:"geometry_id,omitempty"`
	OnlyFinished bool            `cbor:"only_finished,omitempty"`
}

func (q TrainUnloadQuery) StreamKey() string {
	return "train_load_unload_detector_domain.query.train_load_unload.TrainUnloadQuery"
}

func (q TrainUnloadQuery) Handle(repo cqrs.Repository) (any, error) {
	return nil, nil
}

// Зарегистрируем запросы в регистраторе создавая слушателей потоков
func RegisterQueries() *cqrs.CQRSHadler {
	reg := cqrs.GetCQRSHadler()
	reg.Register(&TrainLoadQuery{})
	reg.Register(&TrainUnloadQuery{})

	return reg
}
