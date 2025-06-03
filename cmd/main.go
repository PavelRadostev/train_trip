package main

import (
	"context"

	"github.com/PavelRadostev/train_trip/internal/broker"
	"github.com/PavelRadostev/train_trip/internal/listener"
	"github.com/PavelRadostev/train_trip/internal/redis"
	"github.com/PavelRadostev/train_trip/internal/repository"
	"github.com/PavelRadostev/train_trip/pkg/config"
	"github.com/PavelRadostev/train_trip/pkg/cqrs"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

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
