package broker

import (
	"context"
	"fmt"
	"log"

	// "sync"
	"time"

	"github.com/as-master/train_trip/internal/listener"
	"github.com/as-master/train_trip/pkg/cqrs"
	"github.com/redis/go-redis/v9"
)

type Broker struct {
	redis       *redis.Client
	cqrsHandler *cqrs.CQRSHadler
	// mu          sync.Mutex
	buffer []interface{} // временное хранилище
}

func New(redis *redis.Client, cqrsHandler *cqrs.CQRSHadler) *Broker {
	return &Broker{
		redis:       redis,
		cqrsHandler: cqrsHandler,
		buffer:      make([]interface{}, 0),
	}
}

func (b *Broker) Run(ctx context.Context) {
	const fn = "internal/broker/broker.Run"
	streams := make([]string, 0, len(b.cqrsHandler.GetStreamKeis()))
	ids := make([]string, 0, len(b.cqrsHandler.GetStreamKeis()))
	for _, stream := range b.cqrsHandler.GetStreamKeis() {
		streams = append(streams, stream)
		ids = append(ids, "0") // "0" Прочитать ВСЁ (в т.ч. старое); "$" только новые сообщения, которые появятся после запуска XREAD (блокирующий режим)
	}

	for {
		xres, err := b.redis.XRead(ctx, &redis.XReadArgs{
			Streams: append(streams, ids...),
			Block:   2 * time.Second,
		}).Result()
		if err != nil && err != redis.Nil {
			log.Printf("%s: XRead error: %v", fn, err)
			continue
		}

		for _, stream := range xres {
			for _, msg := range stream.Messages {
				entity, err := b.cqrsHandler.Get(stream.Stream) // создаём новый объект
				if err != nil {
					log.Printf("%s: Unable to clone template for stream %s", fn, stream.Stream)
					continue
				}

				req, err := listener.MsgToReq(msg, entity)
				if err != nil {
					log.Printf("%s: Decode error for stream %s: %v", fn, stream.Stream, err)
					continue
				}
				resp := req.Handle(b.cqrsHandler.GetRepo()) // обработка сообщения
				fmt.Print(resp)

			}
		}
		time.Sleep(100 * time.Microsecond)
	}
}

type Response struct {
	ReqID string
	Data  interface{}
}

type Serde interface {
	Serialize(data interface{}) ([]byte, error)
}

// func (b *Broker) SendResponse(ctx context.Context, resp *Response, timeout time.Duration) error {
// 	// сериализуем ответ
// 	data, err := Serde.Serialize(resp.Data)
// 	if err != nil {
// 		return fmt.Errorf("failed to serialize response: %w", err)
// 	}

// 	// используем pipeline без транзакции
// 	pipe := b.redis.Pipeline()
// 	pipe.RPush(ctx, resp.ReqID, data)
// 	pipe.Expire(ctx, resp.ReqID, timeout)

// 	_, err = pipe.Exec(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute redis pipeline: %w", err)
// 	}
// 	return nil
// }

// processMessage — заглушка: складывает обработанные данные в буфер
// func (b *Broker) processMessage(data *serializer.TransportRequest) {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	b.buffer = append(b.buffer, data)
// 	log.Printf("Buffered message: %#v", data)
// }
