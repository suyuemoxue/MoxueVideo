package wsapi

import (
	"encoding/json"
	"time"
)

type Message struct {
	Type string          `json:"type"`
	ID   string          `json:"id,omitempty"`
	TS   int64           `json:"ts"`
	Data json.RawMessage `json:"data,omitempty"`
}

func NewMessage(t string, data any) Message {
	var raw json.RawMessage
	if data != nil {
		b, _ := json.Marshal(data)
		raw = b
	}
	return Message{Type: t, TS: time.Now().UnixMilli(), Data: raw}
}
