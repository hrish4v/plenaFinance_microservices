package client

import "context"

type UserServiceClient interface {
	CheckAdmin(ctx context.Context, username string) (bool, error)
}
