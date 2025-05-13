package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/as-master/train_trip/internal/config"
	"github.com/as-master/train_trip/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client может быть pgx.Conn или pgxpool.Pool
type Client interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) (pgx.Row, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewCli(cfg config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	const fn = "storage.pg.New"

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
