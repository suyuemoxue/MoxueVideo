package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryUint64(c *gin.Context, key string, defaultValue uint64) uint64 {
	v := c.Query(key)
	if v == "" {
		return defaultValue
	}
	u, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return u
}

func QueryInt(c *gin.Context, key string, defaultValue int, min int, max int) int {
	v := c.Query(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	if i < min {
		return min
	}
	if max > 0 && i > max {
		return max
	}
	return i
}
