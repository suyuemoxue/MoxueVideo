package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func (h *Handler) NotifyWS(c *gin.Context) {
	userID, ok := extractUserID(c.Request)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	websocket.Handler(func(ws *websocket.Conn) {
		_ = ws.SetDeadline(time.Now().Add(70 * time.Second))
		h.notifyHub.add(userID, ws)
		defer func() {
			h.notifyHub.remove(userID, ws)
			_ = ws.Close()
		}()

		_ = websocket.Message.Send(ws, mustJSON(newWSMessage("hello", map[string]any{"user_id": userID})))

		for {
			var raw string
			if err := websocket.Message.Receive(ws, &raw); err != nil {
				return
			}
			_ = ws.SetDeadline(time.Now().Add(70 * time.Second))

			var m WSMessage
			if err := json.Unmarshal([]byte(raw), &m); err != nil {
				continue
			}
			switch m.Type {
			case "ping":
				_ = websocket.Message.Send(ws, mustJSON(newWSMessage("pong", nil)))
			}
		}
	}).ServeHTTP(c.Writer, c.Request)
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
