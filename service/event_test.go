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
	ctx                 context.Context
}

func (suite *EventTestSuite) SetupTest() {
	suite.mockEventRepository = &MockEventRepository{}
	suite.service = NewEvent(suite.mockEventRepository)
	suite.ctx = context.Background()
}

func (suite *EventTestSuite) TestCreateHappyFlow() {
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
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.mockEventRepository.On("Create", suite.ctx, model.Event{UserID: 1, SlotID: 1, InviteeName: "test", InviteeEmail: "test@example.xyz"}).Return(expectedResp, nil)

	resp, err := suite.service.Create(suite.ctx, 1, input)
	suite.Nil(err)
	suite.Equal(1, resp.ID)
	suite.Equal(1, resp.UserID)
	suite.Equal(1, resp.SlotID)
	suite.Equal("test", resp.InviteeName)
}

func (suite *EventTestSuite) TestCreateShouldReturnErrorIfRepositoryReturnsError() {
	input := contract.Event{
		SlotID:       1,
		InviteeName:  "test",
		InviteeEmail: "test@example.xyz",
	}
	suite.mockEventRepository.On("Create", suite.ctx, model.Event{UserID: 1, SlotID: 1, InviteeName: "test", InviteeEmail: "test@example.xyz"}).
		Return(model.Event{}, errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, 1, input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
