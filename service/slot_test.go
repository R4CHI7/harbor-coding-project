package service

import (
	"context"
	"testing"

	"github.com/harbor-xyz/coding-project/model"
	"gorm.io/datatypes"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SlotTestSuite struct {
	suite.Suite
	mockSlotRepository         *MockSlotRepository
	mockAvailabilityRepository *MockUserAvailabilityRepository
	service                    Slot
	ctx                        context.Context
}

func (suite *SlotTestSuite) SetupTest() {
	suite.mockSlotRepository = &MockSlotRepository{}
	suite.mockAvailabilityRepository = &MockUserAvailabilityRepository{}
	suite.service = NewSlot(suite.mockSlotRepository, suite.mockAvailabilityRepository)
	suite.ctx = context.Background()
}

func (suite *SlotTestSuite) TestCreateHappyFlow() {
	suite.mockSlotRepository.On("Get", suite.ctx, 1, mock.Anything, mock.Anything).Return([]model.Slot{}, nil)
	suite.mockAvailabilityRepository.On("Get", suite.ctx, 1).Return(model.UserAvailability{
		UserID: 1,
		Availability: []model.DayAvailability{
			{
				Day:       "monday",
				StartTime: datatypes.NewTime(10, 0, 0, 0),
				EndTime:   datatypes.NewTime(17, 0, 0, 0),
			},
			{
				Day:       "tuesday",
				StartTime: datatypes.NewTime(9, 0, 0, 0),
				EndTime:   datatypes.NewTime(17, 0, 0, 0),
			},
		},
		MeetingDurationMins: 30,
	}, nil)
	suite.mockSlotRepository.On("Create", suite.ctx, mock.Anything).Return(nil)

	numSlots, err := suite.service.CreateSlots(suite.ctx, 1, 14)
	suite.Equal(60, numSlots)
	suite.Nil(err)
}

func TestSlotTestSuite(t *testing.T) {
	suite.Run(t, new(SlotTestSuite))
}
