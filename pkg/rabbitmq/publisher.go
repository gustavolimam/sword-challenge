package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/sword-challenge/pkg/rabbitmq/models"
)

type Publisher interface {
	GetQueueName() (*string, error)
	SendMessage(exchange string, routingKey string, mandatory bool, immediate bool, message models.PublishingMessage) error
}

type publisherImpl struct {
	channel *amqp.Channel
	queue   *amqp.Queue
}

func (queue publisherImpl) SendMessage(exchange string, routingKey string, mandatory bool, immediate bool, message models.PublishingMessage) error {
	if message.ContentType == "" {
		message.ContentType = "application/json"
	}

	return queue.channel.Publish(
		exchange,
		routingKey,
		mandatory,
		immediate,
		amqp.Publishing(message),
	)
}

func (queue publisherImpl) GetQueueName() (*string, error) {
	if queue.queue == nil {
		return nil, fmt.Errorf("not connect to a queue")
	}

	return &queue.queue.Name, nil
}
