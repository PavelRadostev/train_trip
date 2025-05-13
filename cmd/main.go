package main

import (
	"context"

	"github.com/as-master/train_trip/internal/broker"
	"github.com/as-master/train_trip/internal/config"
	"github.com/as-master/train_trip/internal/listener"
	"github.com/as-master/train_trip/internal/redis"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	listeners := listener.RegisterQueries()

	rdb := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	b := broker.New(rdb, listeners)
	b.Run(ctx)

}
