package listener

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type Request struct {
	ReqID        string
	ReturnResult bool
}

func (r *Request) ParseMsg(msg map[string]any) error {
	// 1. request_id
	requestID, ok := msg["i"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'i'")
	}
	r.ReqID = requestID

	// 2. return_result
	returnRaw, ok := msg["r"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'r'")
	}
	r.ReturnResult = returnRaw == "1"

	// 3. main payload
	payload, ok := msg["p"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'p'")
	}

	if err := cbor.Unmarshal([]byte(payload), r); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	return nil
}

type Response struct {
	ReqID string
	Data  interface{}
}
