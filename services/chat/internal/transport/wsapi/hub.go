package wsapi

import (
	"encoding/json"
	"sync"

	"golang.org/x/net/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	conns map[uint64]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{conns: make(map[uint64]map[*websocket.Conn]struct{})}
}

func (h *Hub) Add(userID uint64, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, ok := h.conns[userID]
	if !ok {
		m = make(map[*websocket.Conn]struct{})
		h.conns[userID] = m
	}
	m[c] = struct{}{}
}

func (h *Hub) Remove(userID uint64, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, ok := h.conns[userID]
	if !ok {
		return
	}
	delete(m, c)
	if len(m) == 0 {
		delete(h.conns, userID)
	}
}

func (h *Hub) Send(userID uint64, msg Message) int {
	h.mu.RLock()
	m := h.conns[userID]
	conns := make([]*websocket.Conn, 0, len(m))
	for c := range m {
		conns = append(conns, c)
	}
	h.mu.RUnlock()

	if len(conns) == 0 {
		return 0
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return 0
	}

	n := 0
	for _, c := range conns {
		if err := websocket.Message.Send(c, string(b)); err == nil {
			n++
		}
	}
	return n
}
