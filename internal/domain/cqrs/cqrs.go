package domain

import "fmt"

// Query - базовая пустая структура
type Query[T any] struct{}

// CQRSEntity - интерфейс для CQRS сущностей Query, Command и Event
type CQRSEntity interface {
	StreamKey() string
	Handle(msg string) string
}

// Registrar — регистратор слушателей потоков в Редисе
type Registrar struct {
	entities            map[string]CQRSEntity
	registeredStreamKey []string
}

var registrar = &Registrar{
	entities:            make(map[string]CQRSEntity),
	registeredStreamKey: make([]string, 0),
}

func GetRegistrar() *Registrar {
	return registrar
}

func (r *Registrar) Register(entity CQRSEntity) {
	r.entities[entity.StreamKey()] = entity
	r.registeredStreamKey = append(r.registeredStreamKey, entity.StreamKey())
}

func (r *Registrar) Get(name string) (CQRSEntity, error) {
	entity, ok := r.entities[name]
	if !ok {
		return nil, fmt.Errorf("entity %s not found", name)
	}
	return entity, nil
}

func (r *Registrar) GetStreamKeis() []string {
	return r.registeredStreamKey
}
