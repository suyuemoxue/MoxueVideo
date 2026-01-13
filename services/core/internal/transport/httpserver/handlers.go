package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	deps Dependencies
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

func (h *Handlers) Healthz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbOK := true
	if h.deps.DB != nil {
		sqlDB, err := h.deps.DB.DB()
		if err != nil {
			dbOK = false
		} else if err := sqlDB.PingContext(ctx); err != nil {
			dbOK = false
		}
	}

	redisOK := true
	if h.deps.Redis != nil {
		if err := h.deps.Redis.Ping(ctx).Err(); err != nil {
			redisOK = false
		}
	}

	rabbitOK := h.deps.RabbitMQ != nil && h.deps.RabbitMQ.Conn != nil && !h.deps.RabbitMQ.Conn.IsClosed()

	grpcOK := true
	if h.deps.ChatClient != nil {
		if err := h.deps.ChatClient.Ping(ctx); err != nil {
			grpcOK = false
		}
	}

	status := http.StatusOK
	if !(dbOK && redisOK && rabbitOK && grpcOK) {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"ok":       status == http.StatusOK,
		"mysql":    dbOK,
		"redis":    redisOK,
		"rabbitmq": rabbitOK,
		"grpc":     grpcOK,
	})
}

func (h *Handlers) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) GetUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) ListFollowing(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) ListFollowers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) FollowUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) UnfollowUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) UploadVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) GetVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) LikeVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) UnlikeVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) FavoriteVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) UnfavoriteVideo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}

func (h *Handlers) RecordWatch(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"code": "NOT_IMPLEMENTED"})
}
