package cqrs

import (
	"fmt"

	"github.com/as-master/train_trip/pkg/pgrepo"
)

// CQRSEntity - интерфейс для CQRS сущностей Query, Command и Event
type CQRSEntity interface {
	StreamKey() string
	Handle(repo Repository) (any, error)
}

// Repository - интерфейс для работы с хранилищем данных
type Repository interface {
	GetConnection() pgrepo.Connector
}

// Registrar — регистратор слушателей потоков в Редисе
type CQRSHadler struct {
	entities            map[string]CQRSEntity
	registeredStreamKey []string
	repo                Repository
}

var cqrsHadler = &CQRSHadler{
	entities:            make(map[string]CQRSEntity),
	registeredStreamKey: make([]string, 0),
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
