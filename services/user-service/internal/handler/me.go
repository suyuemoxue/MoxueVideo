package handler

import (
	"net/http"

	"example.com/MoxueVideo/user-service/internal/repo"
	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	users repo.UserRepo
}

func NewMeHandler(users repo.UserRepo) *MeHandler {
	return &MeHandler{users: users}
}

func (h *MeHandler) Me(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}
	userID, ok := v.(uint64)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}

	u, err := h.users.FindByID(c.Request.Context(), userID)
	if err != nil {
		httpx.Fail(c, http.StatusNotFound, httpx.CodeNotFound, "user not found")
		return
	}

	httpx.OK(c, http.StatusOK, u)
}
