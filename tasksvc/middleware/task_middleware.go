package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"tasksvc/config"
	"tasksvc/repository/client"
)

type TaskMiddleware struct {
	userSvcClient client.UserServiceClient
	config        *config.StartupConfig
}

func NewUserMiddleware(userSvcClient client.UserServiceClient, config *config.StartupConfig) *TaskMiddleware {
	return &TaskMiddleware{
		userSvcClient,
		config,
	}
}

func (m *TaskMiddleware) CheckIfAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			//authorizationHeader, _, _ := ctx.Request().BasicAuth()
			//credentialsEncoded := strings.SplitN(authorizationHeader, " ", 2)[1]
			//credentialsDecoded, err := base64.StdEncoding.DecodeString(credentialsEncoded)
			//if err != nil {
			//	fmt.Println(fmt.Sprintf("ERROR :: TaskMiddleware :: CheckIfAdmin :: Error decoding credentials: %v", err.Error()))
			//	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			//}
			//credentials := strings.SplitN(string(credentialsDecoded), ":", 2)
			//username := credentials[0]

			username, _, _ := ctx.Request().BasicAuth()

			fmt.Println("username", username)
			isAdmin, err := m.userSvcClient.CheckAdmin(ctx.Request().Context(), username)
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
