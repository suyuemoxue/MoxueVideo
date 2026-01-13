package mq

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"

	"moxuevideo/chat/internal/domain"
)

type Publisher struct {
	rmq        *RabbitMQ
	exchange   string
	routingKey string
}

func NewPublisher(rmq *RabbitMQ) *Publisher {
	return &Publisher{
		rmq:        rmq,
		exchange:   "moxuevideo.events",
		routingKey: "chat.message.created",
	}
}

func (p *Publisher) PublishChatMessageCreated(evt domain.ChatMessageCreated) error {
	if p == nil || p.rmq == nil || p.rmq.Channel == nil {
		return amqp091.ErrClosed
	}
	if err := p.rmq.Channel.ExchangeDeclare(p.exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return p.rmq.Channel.PublishWithContext(
		nil,
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
}
