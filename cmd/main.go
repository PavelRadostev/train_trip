package main

import (
	"context"
	"log/slog"

	"github.com/PavelRadostev/train_trip/internal/broker"
	"github.com/PavelRadostev/train_trip/internal/listener"
	"github.com/PavelRadostev/train_trip/internal/redis"
	"github.com/PavelRadostev/train_trip/internal/repository"
	"github.com/PavelRadostev/train_trip/pkg/config"
	"github.com/PavelRadostev/train_trip/pkg/cqrs"
	"github.com/PavelRadostev/train_trip/pkg/logger"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	log := logger.SetupLogger("local")
	log.Info("Starting Train Trip Service", slog.String("version", "1.0.0"))

	repo, err := repository.NewPGRepo(cfg, ctx)
	if err != nil {
		panic(err)
	}

	cqrs.InitRepo(repo)
	cqrsHandler := listener.RegisterQueries()

	rdb := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	b := broker.New(rdb, cqrsHandler)
	b.Run(ctx)

}
