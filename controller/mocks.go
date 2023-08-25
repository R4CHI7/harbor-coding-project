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

func (mock *MockUserService) GetAvailabilityOverlap(ctx context.Context, user1ID, user2ID int) (contract.UserAvailabilityOverlap, error) {
	args := mock.Called(ctx, user1ID, user2ID)
	return args.Get(0).(contract.UserAvailabilityOverlap), args.Error(1)
}

type MockEventService struct {
	mock.Mock
}

func (mock *MockEventService) Create(ctx context.Context, userID int, input contract.Event) (contract.EventResponse, error) {
	args := mock.Called(ctx, userID, input)
	return args.Get(0).(contract.EventResponse), args.Error(1)
}

func (mock *MockEventService) GetAll(ctx context.Context, userID int) (contract.EventListResponse, error) {
	args := mock.Called(ctx, userID)
	return args.Get(0).(contract.EventListResponse), args.Error(1)
}

type MockSlotService struct {
	mock.Mock
}

func (mock *MockSlotService) Create(ctx context.Context, userID, numDays int) (int, error) {
	args := mock.Called(ctx, userID, numDays)
	return args.Int(0), args.Error(1)
}

func (mock *MockSlotService) GetAll(ctx context.Context, userID int) (contract.SlotList, error) {
	args := mock.Called(ctx, userID)
	return args.Get(0).(contract.SlotList), args.Error(1)
}

func (mock *MockSlotService) DeleteByID(ctx context.Context, slotID int) error {
	args := mock.Called(ctx, slotID)
	return args.Error(0)
}
