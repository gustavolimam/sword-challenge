package consumeNotify

import (
	"github.com/rs/zerolog/log"
	"github.com/sword-challenge/pkg/rabbitmq/models"
)

type Service interface {
	ConsumeEvent(data models.IncomingEventMessage) bool
}

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) ConsumeEvent(data models.IncomingEventMessage) bool {
	log.Info().Msgf("New task updated: %v", data.Content.Properties)

	return true
}
