package middleware

import (
	"net/http"
	"strings"

	"example.com/MoxueVideo/user-service/internal/service"
	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer"))
		if token == "" {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "missing token")
			c.Abort()
			return
		}

		userID, err := service.ParseUserIDFromToken(token, secret)
		if err != nil {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
