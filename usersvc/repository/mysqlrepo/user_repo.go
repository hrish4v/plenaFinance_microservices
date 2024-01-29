package mysqlrepo

import (
	"context"
	"usersvc/models"
)

//go:generate mockery --name=UserMySQLRepo --inpackage --log-level trace

type UserMySQLRepo interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, username string) (models.User, error)
	DeleteUser(ctx context.Context, username string) error
}
