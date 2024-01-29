package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"usersvc/config"
	"usersvc/service"
)

type UserMiddleware struct {
	userService service.UserService
	config      *config.StartupConfig
}

func NewUserMiddleware(userService service.UserService, config *config.StartupConfig) *UserMiddleware {
	return &UserMiddleware{
		userService,
		config,
	}
}

func (m *UserMiddleware) CheckIfAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			username, _, _ := ctx.Request().BasicAuth()
			fmt.Println("INFO :: UserMiddleware :: CheckIfAdmin :: username:", username)
			isAdmin, err := m.userService.IsAdmin(ctx.Request().Context(), username)
			if err != nil {
				fmt.Println(fmt.Sprintf("ERROR :: UserMiddleware :: CheckIfAdmin :: Error while getting user with username %s: %s", username, err.Error()))
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			if !isAdmin {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c := ctx.Request().Context()
			c = context.WithValue(c, "username", username)

			ctx.SetRequest(ctx.Request().WithContext(c))
			return next(ctx)
		}
	}
}
