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

func (mock *MockUserService) Create(ctx context.Context, input contract.User) (model.User, error) {
	args := mock.Called(ctx, input)
	return args.Get(0).(model.User), args.Error(1)
}
