package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"time"
	"usersvc/config"
	"usersvc/constants"
	httpLocal "usersvc/delivery/http"
	"usersvc/middleware"
	"usersvc/models"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultMysqlDBConnection int = 1
)

func newCfAdminMysqlDbConnection(config *config.StartupConfig) (*models.TaskDB, error) {
	dbCfg := config.TaskDB
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	val := url.Values{}
	val.Add("charset", "utf8")
	val.Add("parseTime", "true")
	val.Add("loc", "Local")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}
	if dbCfg.MaxConnections <= 0 {
		dbCfg.MaxConnections = defaultMysqlDBConnection
	}

	dbConn.SetMaxOpenConns(dbCfg.MaxConnections)
	dbConn.SetMaxIdleConns(dbCfg.MaxConnections)
	dbConn.SetConnMaxLifetime(time.Minute * 3)

	taskDB := &models.TaskDB{}
	taskDB.Db = dbConn.Unsafe()
	fmt.Println("INFO :: Connected to new task db")
	return taskDB, nil
}

func newRouter(UserHandler *httpLocal.UserHandler, userMiddleware *middleware.UserMiddleware) (*echo.Echo, error) {
	e := echo.New()

	internal := e.Group(constants.APIVersionV1)

	e.GET("/actuator/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	internal.POST("/user", UserHandler.AddUser, userMiddleware.CheckIfAdmin())
	internal.DELETE("/user", UserHandler.DeleteUser, userMiddleware.CheckIfAdmin())
	internal.GET("/user", UserHandler.GetUser)
	internal.GET("/check-admin", UserHandler.CheckAdmin, userMiddleware.CheckIfAdmin())
	return e, nil
}

func newApp(config *config.StartupConfig, e *echo.Echo) *App {
	return &App{
		config: config,
		_echo:  e,
	}
}

type App struct {
	config *config.StartupConfig
	_echo  *echo.Echo
}

func (app *App) Start() error {
	return app._echo.Start(app.config.Server.Port)
}
