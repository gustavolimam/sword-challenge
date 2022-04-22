package user

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sword-challenge/internal/model"
	"github.com/sword-challenge/internal/repository/userRepository"
)

type Service interface {
	Login(c echo.Context) error
	GetUserByToken(token string) (*model.User, error)
}

type service struct {
	userQueries userRepository.UserQueries
}

func New(db *sqlx.DB) Service {
	return &service{userRepository.New(db)}
}

func (s *service) Login(c echo.Context) error {
	log.Info().Msg("Start method to login")
	user := &model.User{}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	token, err := s.userQueries.AuthenticateUser(user.ID)
	if err != nil || token == "" {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, token)
}

func (s *service) GetUserByToken(token string) (*model.User, error) {
	return s.userQueries.GetUserByToken(token)
}
