package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"

	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type UserTestSuite struct {
	suite.Suite
	service                        User
	mockUserRepository             *MockUserRepository
	mockUserAvailabilityRepository *MockUserAvailabilityRepository
	ctx                            context.Context
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockUserRepository = &MockUserRepository{}
	suite.mockUserAvailabilityRepository = &MockUserAvailabilityRepository{}
	suite.service = NewUser(suite.mockUserRepository, suite.mockUserAvailabilityRepository)
	suite.ctx = context.Background()
}

func (suite *UserTestSuite) TestCreateHappyFlow() {
	input := contract.User{
		Name:  "test",
		Email: "test@example.xyz",
	}
	expectedResp := model.User{
		ID:        1,
		Name:      "test",
		Email:     "test@example.xyz",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mockUserRepository.On("Create", suite.ctx, model.User{Name: "test", Email: "test@example.xyz"}).Return(expectedResp, nil)

	resp, err := suite.service.Create(suite.ctx, input)
	suite.Nil(err)
	suite.Equal(expectedResp, resp)
}

func (suite *UserTestSuite) TestCreateShouldReturnErrorIfRepositoryFails() {
	input := contract.User{
		Name:  "test",
		Email: "test@example.xyz",
	}
	suite.mockUserRepository.On("Create", suite.ctx, model.User{Name: "test", Email: "test@example.xyz"}).Return(model.User{}, errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, input)
	suite.Error(err, "some error")
	suite.Empty(resp)
}

func (suite *UserTestSuite) TestSetAvailabilityHappyFlow() {
	availability := []model.DayAvailability{
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
	}
	input := contract.UserAvailability{
		Availability:        availability,
		MeetingDurationMins: 30,
	}
	expectedResp := model.UserAvailability{
		ID:                  1,
		UserID:              1,
		Availability:        availability,
		MeetingDurationMins: 30,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	suite.mockUserAvailabilityRepository.On("Set", suite.ctx, model.UserAvailability{
		UserID: 1, Availability: availability, MeetingDurationMins: 30,
	}).Return(expectedResp, nil)

	resp, err := suite.service.SetAvailability(suite.ctx, 1, input)
	suite.Nil(err)
	suite.Equal(expectedResp, resp)
}

func (suite *UserTestSuite) TestSetAvailabilityShouldReturnErrorIfRepositoryFails() {
	availability := []model.DayAvailability{
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
	}
	input := contract.UserAvailability{
		Availability:        availability,
		MeetingDurationMins: 30,
	}

	suite.mockUserAvailabilityRepository.On("Set", suite.ctx, model.UserAvailability{
		UserID: 1, Availability: availability, MeetingDurationMins: 30,
	}).Return(model.UserAvailability{}, errors.New("some error"))

	resp, err := suite.service.SetAvailability(suite.ctx, 1, input)
	suite.Error(err, "some error")
	suite.Empty(resp)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
