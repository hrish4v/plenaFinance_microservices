// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "usersvc/models"

	mock "github.com/stretchr/testify/mock"
)

// MockUserMySQLRepo is an autogenerated mock type for the UserMySQLRepo type
type MockUserMySQLRepo struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *MockUserMySQLRepo) CreateUser(ctx context.Context, user models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: ctx, username
func (_m *MockUserMySQLRepo) DeleteUser(ctx context.Context, username string) error {
	ret := _m.Called(ctx, username)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: ctx, username
func (_m *MockUserMySQLRepo) GetUser(ctx context.Context, username string) (models.User, error) {
	ret := _m.Called(ctx, username)

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockUserMySQLRepo creates a new instance of MockUserMySQLRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserMySQLRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserMySQLRepo {
	mock := &MockUserMySQLRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
