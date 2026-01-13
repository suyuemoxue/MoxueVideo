package httpapi

import (
	"encoding/json"
	"time"
)

type WSMessage struct {
	Type string          `json:"type"`
	ID   string          `json:"id,omitempty"`
	TS   int64           `json:"ts"`
	Data json.RawMessage `json:"data,omitempty"`
}

func newWSMessage(t string, data any) WSMessage {
	var raw json.RawMessage
	if data != nil {
		b, _ := json.Marshal(data)
		raw = b
	}
	return WSMessage{
		Type: t,
		TS:   time.Now().UnixMilli(),
		Data: raw,
	}
}
