package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"moxuevideo/core/internal/usecase/oss"
	"moxuevideo/core/internal/usecase/user"
	"moxuevideo/core/internal/usecase/video"
)

type Checker interface {
	Check(ctx context.Context) error
}

type HealthDeps struct {
	MySQL    Checker
	Redis    Checker
	RabbitMQ Checker
	GRPC     Checker
}

type Handler struct {
	userService  *user.Service
	videoService *video.Service
	ossService   *oss.Service

	health    HealthDeps
	notifyHub *notifyHub
}

type Deps struct {
	User   *user.Service
	Video  *video.Service
	OSS    *oss.Service
	Health HealthDeps
}

func New(deps Deps) *Handler {
	return &Handler{
		userService:  deps.User,
		videoService: deps.Video,
		ossService:   deps.OSS,
		health:       deps.Health,
		notifyHub:    newNotifyHub(),
	}
}

func (h *Handler) Healthz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	mysqlOK := h.health.MySQL == nil || h.health.MySQL.Check(ctx) == nil
	redisOK := h.health.Redis == nil || h.health.Redis.Check(ctx) == nil
	rabbitOK := h.health.RabbitMQ == nil || h.health.RabbitMQ.Check(ctx) == nil
	grpcOK := h.health.GRPC == nil || h.health.GRPC.Check(ctx) == nil

	status := http.StatusOK
	if !(mysqlOK && redisOK && rabbitOK && grpcOK) {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"ok":       status == http.StatusOK,
		"mysql":    mysqlOK,
		"redis":    redisOK,
		"rabbitmq": rabbitOK,
		"grpc":     grpcOK,
	})
}
