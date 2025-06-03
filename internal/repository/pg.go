package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/PavelRadostev/train_trip/pkg/config"
	"github.com/PavelRadostev/train_trip/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	const fn = "internal.repository.pg.NewConn"

	ctx, cancel := context.WithTimeout(ctx, cfg.DB.ConnectTimeout*time.Second)
	defer cancel()

	var pool *pgxpool.Pool
	err := utils.RetriableExecute(func() error {
		var err error
		pool, err = pgxpool.New(ctx, cfg.DB.DSN)
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL [%s]: %w", fn, err)
		}
		// Пинг для проверки соединения
		if err = pool.Ping(ctx); err != nil {
			return fmt.Errorf("ping failed [%s]: %w", fn, err)
		}
		return nil
	}, 5, 1*time.Second)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL [%s]: %w", fn, err)
	}

	return pool, nil
}
