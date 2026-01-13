package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"moxuevideo/core/internal/usecase/oss"
)

func (h *Handler) GetOSSTSToken(c *gin.Context) {
	purpose := c.Query("purpose")
	if purpose == "" {
		purpose = c.Query("biz")
	}
	if purpose == "" {
		purpose = "video"
	}

	userID := uint64(0)
	if v, ok := c.Get("userID"); ok {
		if id, ok := v.(uint64); ok {
			userID = id
		}
	}

	token, err := h.ossService.GetUploadToken(c.Request.Context(), purpose, userID)
	if err != nil {
		if err == oss.ErrUnavailable {
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": "OSS_STS_UNAVAILABLE"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "OSS_STS_ERROR"})
		return
	}
	c.JSON(http.StatusOK, token)
}
