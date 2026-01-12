package middleware

import (
	"net/http"
	"strings"

	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		tokenStr := strings.TrimSpace(strings.TrimPrefix(h, "Bearer"))
		if tokenStr == "" {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "missing token")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		subFloat, ok := sub.(float64)
		if !ok {
			httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("userID", uint64(subFloat))
		c.Next()
	}
}
