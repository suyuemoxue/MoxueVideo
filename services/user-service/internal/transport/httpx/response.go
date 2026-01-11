package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK           = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeNotFound     = 40400
	CodeConflict     = 40900
	CodeInternal     = 50000
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, httpStatus int, data any) {
	c.JSON(httpStatus, Response{Code: CodeOK, Message: "ok", Data: data})
}

func Fail(c *gin.Context, httpStatus int, code int, message string) {
	if httpStatus <= 0 {
		httpStatus = http.StatusInternalServerError
	}
	c.JSON(httpStatus, Response{Code: code, Message: message})
}
