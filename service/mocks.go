package service

import (
	"context"
	"time"

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

type MockEventRepository struct {
	mock.Mock
}

func (mock *MockEventRepository) Create(ctx context.Context, event model.Event) (model.Event, error) {
	args := mock.Called(ctx, event)
	return args.Get(0).(model.Event), args.Error(1)
}

// func (mock *MockEventRepository) Get(ctx context.Context, userID int) ([]model.Event, error) {
// 	args := mock.Called(ctx, userID)
// 	return args.Get(0).([]model.Event), args.Error(1)
// }

type MockSlotRepository struct {
	mock.Mock
}

func (mock *MockSlotRepository) Create(ctx context.Context, slots []model.Slot) error {
	args := mock.Called(ctx, slots)
	return args.Error(0)
}

func (mock *MockSlotRepository) Get(ctx context.Context, userID int, startTimeThreshold, endTimeThreshold time.Time) ([]model.Slot, error) {
	args := mock.Called(ctx, userID, startTimeThreshold, endTimeThreshold)
	return args.Get(0).([]model.Slot), args.Error(1)
}

func (mock *MockSlotRepository) GetByID(ctx context.Context, slotID int) (model.Slot, error) {
	args := mock.Called(ctx, slotID)
	return args.Get(0).(model.Slot), args.Error(1)
}
