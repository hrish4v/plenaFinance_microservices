//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"usersvc/config"
	handler "usersvc/delivery/http"
	middleware "usersvc/middleware"
	mysql_impl "usersvc/repository/mysqlrepo/impl"
	impl "usersvc/service/impl"
)

func InitDependency(config *config.StartupConfig) (*App, error) {
	wire.Build(
		middleware.NewUserMiddleware,
		handler.NewUserHandler,
		impl.NewUserService,
		mysql_impl.NewUserMySQLRepo,
		newRouter,
		newCfAdminMysqlDbConnection,
		newApp)
	return &App{}, nil
}
