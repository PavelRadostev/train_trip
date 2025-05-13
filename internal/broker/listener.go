package broker

// "context"

type Listener interface {
	StreamKey() string
	ParseMsg(map[string]interface{}) (any, error)
	Handle(any) error
	SerializeResp(any)
}
