//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"tasksvc/config"
	handler "tasksvc/delivery/http"
	middleware "tasksvc/middleware"
	client_impl "tasksvc/repository/client/impl"
	mysql_impl "tasksvc/repository/mysqlrepo/impl"
	"tasksvc/service/impl"
	"tasksvc/utils"
)

func InitDependency(config *config.StartupConfig) (*App, error) {
	wire.Build(
		impl.NewTaskService,
		handler.NewTaskHandler,
		client_impl.NewUserServiceClient,
		mysql_impl.NewTaskMySQLRepo,
		utils.NewHttpClient,
		middleware.NewUserMiddleware,
		newRouter,
		newCfAdminMysqlDbConnection,
		newApp)
	return &App{}, nil
}
