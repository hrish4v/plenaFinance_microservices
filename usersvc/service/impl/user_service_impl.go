package impl

import (
	"context"
	"errors"
	"fmt"
	"usersvc/constants"
	"usersvc/models"
	"usersvc/repository/mysqlrepo"
	"usersvc/service"
)

type UserService struct {
	repo mysqlrepo.UserMySQLRepo
}

func NewUserService(repo mysqlrepo.UserMySQLRepo) service.UserService {
	return &UserService{
		repo,
	}
}

func (service UserService) AddUser(ctx context.Context, user models.User) error {
	return service.repo.CreateUser(ctx, user)
}

func (service UserService) DeleteUser(ctx context.Context, username string) error {
	return service.repo.DeleteUser(ctx, username)
}

func (service UserService) GetUser(ctx context.Context, username string) (models.User, error) {
	return service.repo.GetUser(ctx, username)
}

func (service UserService) IsAdmin(ctx context.Context, username string) (bool, error) {
	user, err := service.repo.GetUser(ctx, username)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: UserService :: CheckIfAdmin :: Error while getting user with username %s: %s", username, err.Error()))
		return false, err
	}
	if user.Type == constants.EmptyString {
		fmt.Println(fmt.Sprintf("ERROR :: UserService :: CheckIfAdmin :: No user found with username %s", username))
		return false, errors.New(constants.ErrAccessDenied)
	} else if user.Type == constants.Admin {
		return true, nil
	}
	return false, nil
}
