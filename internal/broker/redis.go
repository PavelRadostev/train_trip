package broker

import (
	"context"
	"log"

	"time"

	"github.com/PavelRadostev/train_trip/pkg/cqrs"
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
		ids = append(ids, "$") // "0" Прочитать ВСЁ (в т.ч. старое); "$" только новые сообщения, которые появятся после запуска XREAD (блокирующий режим)
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
				// Клонируем сущность для обработки
				entity, err := b.cqrsHandler.Get(stream.Stream)
				if err != nil {
					log.Printf("%s: Unable to clone template for stream %s", fn, stream.Stream)
					continue
				}

				// Асинхронная обработка
				go func(msg redis.XMessage, streamName string, entity cqrs.CQRSEntity) {
					reqID, respBytes, ok := b.cqrsHandler.Handle(msg, entity)
					if !ok {
						log.Printf("%s: Failed to handle message from stream %s, msgID=%s", fn, streamName, msg.ID)
						return
					}

					// Ответ в Redis
					pipe := b.redis.Pipeline()
					pipe.RPush(ctx, reqID, respBytes)
					pipe.Expire(ctx, reqID, 30*time.Second)

					// Удаляем сообщение из потока
					pipe.XDel(ctx, streamName, msg.ID)

					if _, err := pipe.Exec(ctx); err != nil {
						log.Printf("%s: Failed to write response or delete message: %v", fn, err)
					}
				}(msg, stream.Stream, entity)
			}
		}
		// time.Sleep(100 * time.Microsecond)
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
