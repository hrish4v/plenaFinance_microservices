package service

import (
	"context"
	"usersvc/models"
)

type UserService interface {
	AddUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, username string) error
	GetUser(ctx context.Context, username string) (models.User, error)
	IsAdmin(ctx context.Context, username string) (bool, error)
}
