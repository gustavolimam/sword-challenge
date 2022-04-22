package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sword-challenge/internal/service/user"
)

func AuthenticatedUser(userService user.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Response().Header().Get(echo.HeaderAuthorization)

			user, err := userService.GetUserByToken(token)
			if err != nil {
				return c.JSON(http.StatusBadRequest, "notAuthenticatedUser")
			}

			c.Set("user", user)

			return next(c)
		}
	}
}

func MustReceiveID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, "idNotReceived")
		}

		c.Set("task_id", id)

		return next(c)
	}
}
