package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/suite"
)

type EventTestSuite struct {
	suite.Suite
	service             Event
	mockEventRepository *MockEventRepository
	mockSlotRepository  *MockSlotRepository
	ctx                 context.Context
}

func (suite *EventTestSuite) SetupTest() {
	suite.mockEventRepository = &MockEventRepository{}
	suite.mockSlotRepository = &MockSlotRepository{}
	suite.service = NewEvent(suite.mockEventRepository, suite.mockSlotRepository)
	suite.ctx = context.Background()
}

func (suite *EventTestSuite) TestCreateHappyFlow() {
	now := time.Now()
	input := contract.Event{
		SlotID:       1,
		InviteeName:  "test",
		InviteeEmail: "test@example.xyz",
	}
	expectedResp := model.Event{
		ID:           1,
		UserID:       1,
		SlotID:       1,
		InviteeName:  "test",
		InviteeEmail: "test@example.xyz",
		StartTime:    now,
		EndTime:      now.Add(30 * time.Minute),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	suite.mockSlotRepository.On("GetByID", suite.ctx, 1).Return(model.Slot{
		ID:        1,
		UserID:    1,
		StartTime: now,
		EndTime:   now.Add(30 * time.Minute),
		Status:    model.StatusCreated,
	}, nil)
	suite.mockEventRepository.On("Create", suite.ctx, model.Event{UserID: 1, SlotID: 1, InviteeName: "test", InviteeEmail: "test@example.xyz", StartTime: now, EndTime: now.Add(30 * time.Minute)}).
		Return(expectedResp, nil)

	resp, err := suite.service.Create(suite.ctx, 1, input)
	suite.Nil(err)
	suite.Equal(1, resp.ID)
	suite.Equal(1, resp.UserID)
	suite.Equal(1, resp.SlotID)
	suite.Equal("test", resp.InviteeName)
}

func (suite *EventTestSuite) TestCreateShouldReturnErrorIfRepositoryReturnsError() {
	now := time.Now()
	input := contract.Event{
		SlotID:       1,
		InviteeName:  "test",
		InviteeEmail: "test@example.xyz",
	}
	suite.mockSlotRepository.On("GetByID", suite.ctx, 1).Return(model.Slot{
		ID:        1,
		UserID:    1,
		StartTime: now,
		EndTime:   now.Add(30 * time.Minute),
		Status:    model.StatusCreated,
	}, nil)
	suite.mockEventRepository.On("Create", suite.ctx, model.Event{UserID: 1, SlotID: 1, InviteeName: "test", InviteeEmail: "test@example.xyz", StartTime: now, EndTime: now.Add(30 * time.Minute)}).
		Return(model.Event{}, errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, 1, input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
