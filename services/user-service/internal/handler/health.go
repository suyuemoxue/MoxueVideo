package handler

import (
	"net/http"

	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	httpx.OK(c, http.StatusOK, gin.H{"status": "ok"})
}
