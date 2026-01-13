package httpserver

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"moxuevideo/core/internal/infra/grpcchat"
	"moxuevideo/core/internal/infra/mq"
)

type Dependencies struct {
	DB         *gorm.DB
	Redis      *redis.Client
	RabbitMQ   *mq.RabbitMQ
	ChatClient *grpcchat.Client
}
