package api

import (
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	cValidator "github.com/sword-challenge/internal/api/validator"
	"github.com/sword-challenge/internal/service/task"
	"github.com/sword-challenge/internal/service/user"
	"github.com/sword-challenge/pkg/rabbitmq"
)

type Routers interface {
	RegisterRoutes()
}

type Router struct {
	base  *echo.Echo
	userS user.Service
	taskS task.Service
}

func Start(e *echo.Echo, db *sqlx.DB, publisher rabbitmq.Publisher, consumer rabbitmq.Consumer) Routers {
	e.Validator = &cValidator.CustomValidator{Validator: validator.New()}
	userS := user.New(db)
	taskS := task.New(db, publisher)

	return Router{e, userS, taskS}
}
