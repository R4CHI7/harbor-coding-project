package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (mock *MockUserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	args := mock.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

type MockUserAvailabilityRepository struct {
	mock.Mock
}

func (mock *MockUserAvailabilityRepository) Set(ctx context.Context, ua model.UserAvailability) (model.UserAvailability, error) {
	args := mock.Called(ctx, ua)
	return args.Get(0).(model.UserAvailability), args.Error(1)
}

func (mock *MockUserAvailabilityRepository) Get(ctx context.Context, userID int) (model.UserAvailability, error) {
	args := mock.Called(ctx, userID)
	return args.Get(0).(model.UserAvailability), args.Error(1)
}
