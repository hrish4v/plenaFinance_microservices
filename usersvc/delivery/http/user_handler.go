package http

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
	"usersvc/constants"
	"usersvc/models"
	"usersvc/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service,
	}
}

func (h *UserHandler) AddUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	userReq := models.User{}
	if err = c.Bind(&userReq); err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to parse the request :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	err = h.service.AddUser(ctx, userReq)
	if err != nil {
		errMsg := err.Error()
		var mysqlErr *mysql.MySQLError
		fmt.Println(fmt.Sprintf("ERROR :: Failed to add the user :: %s", err.Error()))
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			errMsg = constants.ErrDuplicateEmailOrUsername
		}
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: errMsg})
	}
	return c.JSON(http.StatusCreated, &models.StatusMessageResponse{Status: "success", Message: "User added successfully"})
}

func (h *UserHandler) DeleteUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	username := c.Request().URL.Query().Get("username")
	err = h.service.DeleteUser(ctx, username)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to delete the user :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Message: "user deleted successfully", Status: "success"})
}

func (h *UserHandler) GetUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	username := c.Request().URL.Query().Get("username")
	user, err := h.service.GetUser(ctx, username)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to fetch the user :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Data: user, Status: "success"})
}

func (h *UserHandler) CheckAdmin(c echo.Context) (err error) {
	ctx := c.Request().Context()
	username := ctx.Value("username")
	if username == constants.EmptyString {
		fmt.Println("ERROR :: username is not present")
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: constants.ErrUsernameNotPresent})
	}
	adminResponse := models.AdminResponse{IsAdmin: true}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Data: adminResponse, Status: "success"})
}
