package wsapi

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"moxuevideo/chat/internal/usecase/chat"
)

type Server struct {
	hub     *Hub
	service *chat.Service
}

type SendPayload struct {
	ToUserID uint64 `json:"to_user_id"`
	MsgType  string `json:"msg_type"`
	Content  string `json:"content"`
	Uniqued  string `json:"uniqued"`
}

func NewServer(hub *Hub, svc *chat.Service) *Server {
	if hub == nil {
		hub = NewHub()
	}
	return &Server{hub: hub, service: svc}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/ws/chat", websocket.Handler(s.chatWS))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	return mux
}

func (s *Server) chatWS(ws *websocket.Conn) {
	userID, ok := ExtractUserID(ws.Request())
	if !ok {
		_ = ws.Close()
		return
	}

	_ = ws.SetDeadline(time.Now().Add(70 * time.Second))
	s.hub.Add(userID, ws)
	defer func() {
		s.hub.Remove(userID, ws)
		_ = ws.Close()
	}()

	_ = websocket.Message.Send(ws, mustJSON(NewMessage("hello", map[string]any{"user_id": userID})))

	for {
		var raw string
		if err := websocket.Message.Receive(ws, &raw); err != nil {
			return
		}
		_ = ws.SetDeadline(time.Now().Add(70 * time.Second))

		var m Message
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			continue
		}
		switch m.Type {
		case "ping":
			_ = websocket.Message.Send(ws, mustJSON(NewMessage("pong", nil)))
		case "send":
			var p SendPayload
			if err := json.Unmarshal(m.Data, &p); err != nil {
				continue
			}
			if s.service == nil {
				continue
			}
			evt, err := s.service.Send(ws.Request().Context(), userID, p.ToUserID, p.MsgType, p.Content, p.Uniqued)
			if err != nil {
				continue
			}
			_ = websocket.Message.Send(ws, mustJSON(NewMessage("ack", map[string]any{"id": m.ID, "message_id": evt.MessageID, "uniqued": evt.Uniqued})))
			msg := NewMessage("message", map[string]any{
				"message_id": evt.MessageID,
				"sender_id":  evt.SenderID,
				"to_user_id": evt.ReceiverID,
				"msg_type":   evt.MsgType,
				"content":    evt.Content,
				"uniqued":    evt.Uniqued,
				"created_at": evt.CreatedAt,
			})
			s.hub.Send(evt.SenderID, msg)
			if evt.ReceiverID != 0 && evt.ReceiverID != evt.SenderID {
				s.hub.Send(evt.ReceiverID, msg)
			}
		}
	}
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
