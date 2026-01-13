package health

import (
	"context"
	"errors"

	"moxuevideo/core/internal/infra/mq"
)

type RabbitMQChecker struct {
	rmq *mq.RabbitMQ
}

func NewRabbitMQChecker(rmq *mq.RabbitMQ) *RabbitMQChecker {
	return &RabbitMQChecker{rmq: rmq}
}

func (c *RabbitMQChecker) Check(_ context.Context) error {
	if c == nil || c.rmq == nil || c.rmq.Conn == nil {
		return errors.New("rabbitmq unavailable")
	}
	if c.rmq.Conn.IsClosed() {
		return errors.New("rabbitmq connection closed")
	}
	return nil
}
