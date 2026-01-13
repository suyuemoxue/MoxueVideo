package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("X-User-ID")
		if raw == "" {
			raw = c.GetHeader("Authorization")
			raw = strings.TrimSpace(strings.TrimPrefix(raw, "Bearer"))
		}
		raw = strings.TrimSpace(raw)
		if raw == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED"})
			c.Abort()
			return
		}
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED"})
			c.Abort()
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}
