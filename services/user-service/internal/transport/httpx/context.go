package httpx

import "github.com/gin-gonic/gin"

func UserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}
