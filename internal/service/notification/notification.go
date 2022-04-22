package notification

import (
	"encoding/json"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sword-challenge/internal/environment"
	"github.com/sword-challenge/internal/model"
	"github.com/sword-challenge/pkg/rabbitmq"
	"github.com/sword-challenge/pkg/rabbitmq/models"
)

type Service interface {
	TaskNotification(task model.Notification) error
}

type service struct {
	publisher rabbitmq.Publisher
}

func New(publisher rabbitmq.Publisher) Service {
	return &service{publisher}
}

func (s *service) TaskNotification(task model.Notification) error {
	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Warn().Msgf("Failed to marshal task to JSON when sending notification error", err)
		return err
	}

	if err := s.publisher.SendMessage("", os.Getenv(environment.TaskQueue), false, false,
		models.PublishingMessage{ContentType: echo.MIMEApplicationJSON, Body: taskJson}); err != nil {
		log.Warn().Msgf("Failed to publish task completion notification error", err)
		return err
	}

	return nil
}
