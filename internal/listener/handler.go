package listener

import (
	"fmt"

	"github.com/as-master/train_trip/pkg/cqrs"
	"github.com/fxamacker/cbor/v2"
	"github.com/redis/go-redis/v9"
)

type Request struct {
	ReqID        string
	ReturnResult bool
	CQRSEntity   cqrs.CQRSEntity
}

type Response struct {
	ReqID string
	Data  interface{}
}

// Формируем CQRS запрос
func MsgToReq(redisMsg redis.XMessage, entity cqrs.CQRSEntity) (Request, error) {
	const fn = "internal.listener.serializer.MsgToReq"
	msg := redisMsg.Values
	if len(msg) == 0 {
		return Request{}, fmt.Errorf("%s: empty message", fn)
	}
	req := Request{CQRSEntity: entity}
	// 1. request_id
	requestID, ok := msg["i"].(string)
	if !ok {
		return req, fmt.Errorf("%s: missing or invalid 'i'", fn)
	}
	req.ReqID = requestID

	// 2. return_result
	returnRaw, ok := msg["r"].(string)
	if !ok {
		return req, fmt.Errorf("%s: missing or invalid 'r'", fn)
	}
	req.ReturnResult = returnRaw == "1"

	// 3. main payload
	payload, ok := msg["p"].(string)
	if !ok {
		return Request{}, fmt.Errorf("%s: missing or invalid 'p'", fn)
	}

	if err := cbor.Unmarshal([]byte(payload), &entity); err != nil {
		return req, fmt.Errorf("%s: failed to decode payload: %w", fn, err)
	}

	return req, nil
}

func (r Request) Handle(repo cqrs.Repository) Response {
	result, err := r.CQRSEntity.Handle(repo)
	if err != nil {
		return Response{ReqID: r.ReqID, Data: err}
	}
	fmt.Println(result)
	return Response{}
}
