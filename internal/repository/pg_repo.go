package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/as-master/train_trip/internal/domain"
	"github.com/as-master/train_trip/internal/domain/model"
	"github.com/as-master/train_trip/pkg/config"
	"github.com/as-master/train_trip/pkg/pgrepo"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repo struct {
	schema string
	conn   pgrepo.Connector
	logger *log.Logger
}

func (r *Repo) GetConnection() pgrepo.Connector {
	return r.conn
}

func (r *Repo) GetBDSchema() string {
	return r.schema
}

func NewPGRepo(cfg *config.Config, ctx context.Context) (*Repo, error) {
	conn, err := NewPool(cfg, ctx)
	if err != nil {
		panic(fmt.Errorf("failed to create pgx pool: %w", err))
	}
	return &Repo{conn: conn, schema: cfg.DB.Schema}, nil
}

func WrapPGError(err error, logger *log.Logger) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		wrapped := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
		logger.Println(wrapped)
		return wrapped
	}
	return err
}

// LoadByTimeInterval возвращает писок погрузок по интервалу времени и списку идентификаторов поездов
func (r *Repo) LoadByTimeInterval(ctx context.Context, trainIDs []int, timeInterval domain.Interval) ([]model.TrainLoad, error) {
	query := `SELECT 
		id, is_deleted, load_arrive_time, load_begin_time, load_end_time,
		load_depart_time, train_id, geometry_id, unload_id, shovel_id,
		manual, load_type_id_manual, work_type_id_manual,
		load_begin_time_manual, load_end_time_manual,
		unload_id_manual, shovel_id_manual, volume_manual,
		cycle_ids, is_cured, carriage_num, source
	FROM train_load_unload_store.loads 
	WHERE train_id = ANY($1) 
	  AND load_begin_time >= $2 
	  AND load_end_time <= $3`

	rows, err := r.conn.Query(ctx, query, trainIDs, timeInterval.TimeFrom, timeInterval.TimeTo)
	if err != nil {
		return nil, WrapPGError(err, r.logger)
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
			return nil, WrapPGError(err, r.logger)
		}
		result = append(result, load)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapPGError(err, r.logger)
	}

	return result, nil

}

// LoadByID возвращает погрузку по id
func (r *Repo) LoadByID(ctx context.Context, id int) (*model.TrainLoad, error) {
	query := `SELECT 
		id, is_deleted, load_arrive_time, load_begin_time, load_end_time,
		load_depart_time, train_id, geometry_id, unload_id, shovel_id,
		manual, load_type_id_manual, work_type_id_manual,
		load_begin_time_manual, load_end_time_manual,
		unload_id_manual, shovel_id_manual, volume_manual,
		cycle_ids, is_cured, carriage_num, source
		FROM train_load_unload_store.loads WHERE id = $1`

	row := r.conn.QueryRow(ctx, query, id)
	var load model.TrainLoad
	if err := row.Scan(
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
		return nil, WrapPGError(err, r.logger)
	}
	return &load, nil

}

// Update обновляет существующую запись о погрузке
func (r *Repo) Update(ctx context.Context, load *model.TrainLoad) error {
	query := `
		UPDATE train_load_unload_store.loads SET
			is_deleted = $1,
			load_arrive_time = $2,
			load_begin_time = $3,
			load_end_time = $4,
			load_depart_time = $5,
			train_id = $6,
			geometry_id = $7,
			unload_id = $8,
			shovel_id = $9,
			manual = $10,
			load_type_id_manual = $11,
			work_type_id_manual = $12,
			load_begin_time_manual = $13,
			load_end_time_manual = $14,
			unload_id_manual = $15,
			shovel_id_manual = $16,
			volume_manual = $17,
			cycle_ids = $18,
			is_cured = $19,
			carriage_num = $20,
			source = $21
		WHERE id = $22
	`

	tag, err := r.conn.Exec(ctx, query,
		load.IsDeleted,
		load.LoadArriveTime,
		load.LoadBeginTime,
		load.LoadEndTime,
		load.LoadDepartTime,
		load.TrainID,
		load.GeometryID,
		load.UnloadID,
		load.ShovelID,
		load.Manual,
		load.LoadTypeIDManual,
		load.WorkTypeIDManual,
		load.LoadBeginTimeManual,
		load.LoadEndTimeManual,
		load.UnloadIDManual,
		load.ShovelIDManual,
		load.VolumeManual,
		load.CycleIDs,
		load.IsCured,
		load.CarriageNum,
		load.Source,
		load.ID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no TrainLoad found with id %d", load.ID)
	}

	return nil
}

// Delete удаляет запись логически (is_deleted = true)
func (r *Repo) Delete(ctx context.Context, id int) error {
	query := `UPDATE train_load_unload_store.loads SET is_deleted = true WHERE id = $1`

	tag, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no TrainLoad found with id %d", id)
	}

	return nil
}

func (r *Repo) GetByTimeInterva(ctx context.Context, id int) (*model.TrainLoad, error) {
	// Placeholder implementation
	return nil, fmt.Errorf("not implemented")
}
