package controller

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mock *MockUserService) Create(ctx context.Context, input contract.User) (contract.UserResponse, error) {
	args := mock.Called(ctx, input)
	return args.Get(0).(contract.UserResponse), args.Error(1)
}

func (mock *MockUserService) SetAvailability(ctx context.Context, userID int, input contract.UserAvailability) (model.UserAvailability, error) {
	args := mock.Called(ctx, userID, input)
	return args.Get(0).(model.UserAvailability), args.Error(1)
}

func (mock *MockUserService) GetAvailability(ctx context.Context, userID int) (contract.UserAvailability, error) {
	args := mock.Called(ctx, userID)
	return args.Get(0).(contract.UserAvailability), args.Error(1)
}
