package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UploadVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) GetVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) LikeVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) UnlikeVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) FavoriteVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) UnfavoriteVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) CommentVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handler) RecordWatch(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}
