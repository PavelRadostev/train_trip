package broker

import (
	"context"
	"log"
	"sync"
	"time"

	domain "github.com/as-master/train_trip/internal/domain/cqrs"
	"github.com/as-master/train_trip/internal/serializer"
	"github.com/redis/go-redis/v9"
)

type Broker struct {
	redis     *redis.Client
	listeners *domain.Registrar
	mu        sync.Mutex
	buffer    []interface{} // временное хранилище
}

func New(redis *redis.Client, listeners *domain.Registrar) *Broker {
	return &Broker{
		redis:     redis,
		listeners: listeners,
		buffer:    make([]interface{}, 0),
	}
}

func (b *Broker) Run(ctx context.Context) {
	const fn = "internal/broker/broker.Run"
	streams := make([]string, 0, len(b.listeners.GetStreamKeis()))
	ids := make([]string, 0, len(b.listeners.GetStreamKeis()))
	for _, stream := range b.listeners.GetStreamKeis() {
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
				entity, err := b.listeners.Get(stream.Stream) // создаём новый объект
				if err != nil {
					log.Printf("%s: Unable to clone template for stream %s", fn, stream.Stream)
					continue
				}

				transportReq, err := serializer.RadisMsgToTransportReq(msg, entity)
				if err != nil {
					log.Printf("%s: Decode error for stream %s: %v", fn, stream.Stream, err)
					continue
				}

				b.processMessage(transportReq)
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
func (b *Broker) processMessage(data *serializer.TransportRequest) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buffer = append(b.buffer, data)
	log.Printf("Buffered message: %#v", data)
}
