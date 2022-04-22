package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/sword-challenge/pkg/rabbitmq/models"
	"golang.org/x/sync/errgroup"
)

type Consumer interface {
	SubscribeEvents(ctx context.Context, consumerEvent models.ConsumerEvent) error
	ReadMessage(ctx context.Context, correlationID string, consumerEvent models.ConsumerEvent) error
}

type consumerImpl struct {
	queueName string
	channel   *amqp.Channel
}

func (consumer consumerImpl) SubscribeEvents(ctx context.Context, consumerEvent models.ConsumerEvent) error {
	errs, errCtx := errgroup.WithContext(ctx)
	errs.Go(func() error {
		return consumer.getEvents(errCtx, consumerEvent)
	})

	return errs.Wait()
}

func (consumer consumerImpl) getEvents(ctx context.Context, consumerEvent models.ConsumerEvent) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			messages, err := consumer.channel.Consume(consumer.queueName, "", false, false, false, false, nil)
			if err != nil {
				return err
			}

			for message := range messages {
				var body []byte
				if message.Body != nil {
					body = message.Body
				}
				var event models.IncomingEventMessage
				if err := json.Unmarshal(body, &event); err != nil {
					message.Nack(false, false) //To move to dlq we need to send a Nack with requeue = false
					continue
				}

				success := consumerEvent.Handler(event)
				if success {
					message.Ack(true)
				} else {
					message.Nack(false, false) //To move to dlq we need to send a Nack with requeue = false
				}
			}
		}
	}
}

func (consumer consumerImpl) ReadMessage(ctx context.Context, correlationID string, consumerEvent models.ConsumerEvent) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			messages, err := consumer.channel.Consume(consumer.queueName, "", false, false, false, false, nil)
			if err != nil {
				return err
			}

			for message := range messages {
				var body []byte
				if message.Body != nil {
					body = message.Body
				}
				var event models.IncomingEventMessage
				if err := json.Unmarshal(body, &event); err != nil {
					message.Nack(false, false) //To move to dlq we need to send a Nack with requeue = false
					continue
				}

				if message.CorrelationId == correlationID {
					success := consumerEvent.Handler(event)

					if success {
						message.Ack(true)
						return nil
					} else {
						message.Nack(false, true) //Requeue true because we are searching by specif message
					}
				} else {
					message.Nack(false, true) //Requeue true because we are searching by specif message
				}
			}
		}
	}
}
