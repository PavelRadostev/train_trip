package cqrs

import (
	"fmt"

	"github.com/PavelRadostev/train_trip/pkg/pgrepo"
	"github.com/fxamacker/cbor/v2"
	"github.com/redis/go-redis/v9"
)

// CQRSEntity - интерфейс для CQRS сущностей Query, Command и Event
type CQRSEntity interface {
	StreamKey() string
	Handle(repo Repository) (any, error)
}

// Repository - интерфейс для работы с хранилищем данных
type Repository interface {
	GetBDSchema() string
	GetConnection() pgrepo.Connector
}

type Serialiser interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type cborSerialiser struct {
	enc cbor.EncMode
	dec cbor.DecMode
}

func (s *cborSerialiser) Marshal(v any) ([]byte, error) {
	return s.enc.Marshal(v)
}

func (s *cborSerialiser) Unmarshal(data []byte, v any) error {
	return s.dec.Unmarshal(data, v)
}

func tuneSerialiser() Serialiser {
	encMode, err := cbor.EncOptions{
		TimeTag: cbor.EncTagRequired,
	}.EncMode()
	if err != nil {
		panic(fmt.Errorf("failed to create CBOR EncMode: %v", err))
	}

	decMode, err := cbor.DecOptions{}.DecMode()
	if err != nil {
		panic(fmt.Errorf("failed to create CBOR DecMode: %v", err))
	}

	return &cborSerialiser{
		enc: encMode,
		dec: decMode,
	}
}

type Request struct {
	ReqID        string
	ReturnResult bool
	CQRSEntity   CQRSEntity
}

type Response struct {
	ReqID      string `cbor:"req_id"`
	Result     any    `cbor:"result"`
	RespErr    string `cbor:"error,omitempty"`
	RespErrCls string `cbor:"error_class,omitempty"`
}

// CQRSHadler — обработчик зарегистрированных слушателей потоков в Редисе
type CQRSHadler struct {
	entities            map[string]CQRSEntity
	registeredStreamKey []string
	repo                Repository
	serializer          Serialiser
}

var cqrsHadler = &CQRSHadler{
	entities:            make(map[string]CQRSEntity),
	registeredStreamKey: make([]string, 0),
	serializer:          tuneSerialiser(),
}

func GetCQRSHadler() *CQRSHadler {
	return cqrsHadler
}

func InitRepo(dbConn Repository) {
	cqrsHadler.repo = dbConn
}

func (h *CQRSHadler) Register(entity CQRSEntity) {
	h.entities[entity.StreamKey()] = entity
	h.registeredStreamKey = append(h.registeredStreamKey, entity.StreamKey())
}

func (h *CQRSHadler) Get(name string) (CQRSEntity, error) {
	entity, ok := h.entities[name]
	if !ok {
		return nil, fmt.Errorf("entity %s not found", name)
	}
	return entity, nil
}

func (h *CQRSHadler) GetStreamKeis() []string {
	return h.registeredStreamKey
}

func (h *CQRSHadler) GetRepo() Repository {
	return h.repo
}

// Формируем CQRS запрос
func (h CQRSHadler) MsgToReq(redisMsg redis.XMessage, entity CQRSEntity) (Request, error) {
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

	if err := h.serializer.Unmarshal([]byte(payload), &entity); err != nil {
		return req, fmt.Errorf("%s: failed to decode payload: %w", fn, err)
	}

	return req, nil
}

func (r Request) Handle(repo Repository) Response {

	result, err := r.CQRSEntity.Handle(repo)
	if err != nil {
		return Response{ReqID: r.ReqID, Result: err}
	}
	return Response{ReqID: r.ReqID, Result: result}
}

func (h CQRSHadler) Handle(redisMsg redis.XMessage, entity CQRSEntity) (reqID string, result []byte, ok bool) {
	op := "pkg.cqrs.CQRSHadler.Handle"
	ok = true

	var response any

	defer func() {
		data, err := h.serializer.Marshal(response)
		if err != nil {
			ok = false
			// На случай сбоя сериализатора, возвращаем raw CBOR строки
			fallback, _ := cbor.Marshal(fmt.Sprintf("%s: failed to serialize response: %v", op, err))
			result = fallback
			return
		}
		result = data
	}()

	req, err := h.MsgToReq(redisMsg, entity)
	if err != nil {
		ok = false
		response = fmt.Sprintf("%s: failed to decode message: %v", op, err)
		return
	}
	reqID = req.ReqID

	resp, err := req.CQRSEntity.Handle(h.GetRepo())
	if err != nil {
		ok = false
		response = fmt.Sprintf("%s: failed to handle message: %v", op, err)
		return
	}

	response = Response{
		ReqID:  req.ReqID,
		Result: resp,
	}
	return
}
