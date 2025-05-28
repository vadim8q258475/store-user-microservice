package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vadim8q258475/store-user-microservice/config"
)

type rabbitmqPublisher struct {
	queueName string
	channel   *amqp.Channel
}

func NewRabbitMQQueue(channel *amqp.Channel, cfg config.Config) Publisher {
	return &rabbitmqPublisher{
		queueName: cfg.RabbitMQQueueName,
		channel:   channel,
	}
}

func (q *rabbitmqPublisher) Publish(ctx context.Context, data []byte) error {
	return q.channel.PublishWithContext(
		ctx,
		"",
		q.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		},
	)
}
