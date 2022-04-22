package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/sword-challenge/pkg/rabbitmq/models"
)

type Client interface {
	NewPublisher(args *models.QueueArgs) (Publisher, error)
	NewConsumer(queueName string) (Consumer, error)
}

type clientImpl struct {
	connection *amqp.Connection
}

func New(credential models.Credential) (Client, error) {
	conn, err := amqp.Dial(credential.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}

	return clientImpl{connection: conn}, nil
}

func (client clientImpl) NewPublisher(args *models.QueueArgs) (Publisher, error) {
	if args == nil {
		args = &models.QueueArgs{}
	}

	channel, err := client.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel connection error: %v", err)
	}

	queue, err := channel.QueueDeclare(args.Name, args.Durable, args.AutoDelete, args.Exclusive, args.NoWait, nil)
	if err != nil {
		return nil, fmt.Errorf("queue connection error: %v", err)
	}

	return publisherImpl{channel: channel, queue: &queue}, nil
}

func (client clientImpl) NewConsumer(queueName string) (Consumer, error) {
	channel, err := client.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel connection error: %v", err)
	}

	return consumerImpl{channel: channel, queueName: queueName}, nil
}
