package handler

import (
	"errors"
	"net/http"

	"example.com/MoxueVideo/user-service/internal/service"
	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid request")
		return
	}

	userID, err := h.auth.Register(c.Request.Context(), service.RegisterInput{Username: req.Username, Password: req.Password})
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			httpx.Fail(c, http.StatusConflict, httpx.CodeConflict, "username taken")
			return
		}
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	httpx.OK(c, http.StatusCreated, gin.H{"userId": userID})
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid request")
		return
	}

	res, err := h.auth.Login(c.Request.Context(), service.LoginInput{Username: req.Username, Password: req.Password})
	if err != nil {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "invalid credentials")
		return
	}

	httpx.OK(c, http.StatusOK, gin.H{"token": res.Token, "user": res.User})
}
