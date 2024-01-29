package impl

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"usersvc/models"
	"usersvc/repository/mysqlrepo/mocks"
)

var userRepo *mocks.MockUserMySQLRepo
var serviceMock UserService
var ctxMock context.Context

func init() {
	userRepo = new(mocks.MockUserMySQLRepo)

	serviceMock = UserService{
		userRepo,
	}
	ctxMock = context.Background()
}

func Test_AddUser(t *testing.T) {
	// Test 1
	userRepo.On("CreateUser", ctxMock, models.User{}).Return(nil)
	err := serviceMock.AddUser(ctxMock, models.User{})
	assert.NoError(t, err)

	// Test 2
	userRepo.On("CreateUser", ctxMock, models.User{Email: "test@test.com"}).Return(errors.New("error"))
	err = serviceMock.AddUser(ctxMock, models.User{Email: "test@test.com"})
	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
}
