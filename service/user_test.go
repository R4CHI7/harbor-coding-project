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
	suite.Equal(contract.UserResponse{ID: 1}, resp)
}

func (suite *UserTestSuite) TestCreateShouldReturnErrorIfRepositoryFails() {
	input := contract.User{
		Name:  "test",
		Email: "test@example.xyz",
	}
	suite.mockUserRepository.On("Create", suite.ctx, model.User{Name: "test", Email: "test@example.xyz"}).Return(model.User{}, errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, input)
	suite.Equal("some error", err.Error())
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
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *UserTestSuite) TestGetAvailability() {
	availability := model.UserAvailability{
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
	}

	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 1).Return(availability, nil)

	resp, err := suite.service.GetAvailability(suite.ctx, 1)

	suite.Nil(err)
	suite.Equal(availability.MeetingDurationMins, resp.MeetingDurationMins)
}

func (suite *UserTestSuite) TestGetAvailabilityReturnsErrorWhenRepositoryReturnsError() {
	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 1).Return(model.UserAvailability{}, errors.New("some error"))

	resp, err := suite.service.GetAvailability(suite.ctx, 1)

	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *UserTestSuite) TestGerAvailabilityOverlapReturnsOverlapIfItExists() {
	availability1 := []model.DayAvailability{
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
		{
			Day:       "thursday",
			StartTime: datatypes.NewTime(8, 0, 0, 0),
			EndTime:   datatypes.NewTime(20, 0, 0, 0),
		},
	}
	availability2 := []model.DayAvailability{
		{
			Day:       "monday",
			StartTime: datatypes.NewTime(11, 15, 0, 0),
			EndTime:   datatypes.NewTime(18, 0, 0, 0),
		},
		{
			Day:       "wednesday",
			StartTime: datatypes.NewTime(11, 0, 0, 0),
			EndTime:   datatypes.NewTime(17, 0, 0, 0),
		},
		{
			Day:       "thursday",
			StartTime: datatypes.NewTime(9, 0, 0, 0),
			EndTime:   datatypes.NewTime(18, 30, 0, 0),
		},
	}

	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 1).Return(model.UserAvailability{UserID: 1, Availability: availability1, MeetingDurationMins: 30}, nil)
	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 2).Return(model.UserAvailability{UserID: 2, Availability: availability2, MeetingDurationMins: 30}, nil)

	expectedResp := contract.UserAvailabilityOverlap{
		Overlap: []model.DayAvailability{
			{
				Day:       "monday",
				StartTime: datatypes.NewTime(11, 15, 0, 0),
				EndTime:   datatypes.NewTime(17, 0, 0, 0),
			},
			{
				Day:       "thursday",
				StartTime: datatypes.NewTime(9, 0, 0, 0),
				EndTime:   datatypes.NewTime(18, 30, 0, 0),
			},
		},
	}
	resp, err := suite.service.GetAvailabilityOverlap(suite.ctx, 1, 2)
	suite.Nil(err)
	suite.Equal(2, len(resp.Overlap))
	suite.Equal(expectedResp, resp, "This is an undeterministic test, please rerun") // TODO: This is a bit undeterministic due to order of the slice.
}

func (suite *UserTestSuite) TestGerAvailabilityOverlapReturnsNoOverlapIfItDoesNotExist() {
	availability1 := []model.DayAvailability{
		{
			Day:       "monday",
			StartTime: datatypes.NewTime(10, 0, 0, 0),
			EndTime:   datatypes.NewTime(13, 0, 0, 0),
		},
		{
			Day:       "tuesday",
			StartTime: datatypes.NewTime(9, 0, 0, 0),
			EndTime:   datatypes.NewTime(17, 0, 0, 0),
		},
		{
			Day:       "thursday",
			StartTime: datatypes.NewTime(18, 0, 0, 0),
			EndTime:   datatypes.NewTime(23, 0, 0, 0),
		},
	}
	availability2 := []model.DayAvailability{
		{
			Day:       "monday",
			StartTime: datatypes.NewTime(14, 15, 0, 0),
			EndTime:   datatypes.NewTime(18, 0, 0, 0),
		},
		{
			Day:       "wednesday",
			StartTime: datatypes.NewTime(11, 0, 0, 0),
			EndTime:   datatypes.NewTime(17, 0, 0, 0),
		},
		{
			Day:       "thursday",
			StartTime: datatypes.NewTime(9, 0, 0, 0),
			EndTime:   datatypes.NewTime(17, 30, 0, 0),
		},
	}

	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 1).Return(model.UserAvailability{UserID: 1, Availability: availability1, MeetingDurationMins: 30}, nil)
	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 2).Return(model.UserAvailability{UserID: 2, Availability: availability2, MeetingDurationMins: 30}, nil)

	resp, err := suite.service.GetAvailabilityOverlap(suite.ctx, 1, 2)
	suite.Nil(err)
	suite.Equal(0, len(resp.Overlap))
}

func (suite *UserTestSuite) TestGerAvailabilityOverlapReturnsErrorIfRepositoryReturnsError() {
	availability1 := []model.DayAvailability{
		{
			Day:       "monday",
			StartTime: datatypes.NewTime(10, 0, 0, 0),
			EndTime:   datatypes.NewTime(13, 0, 0, 0),
		},
		{
			Day:       "tuesday",
			StartTime: datatypes.NewTime(9, 0, 0, 0),
			EndTime:   datatypes.NewTime(17, 0, 0, 0),
		},
		{
			Day:       "thursday",
			StartTime: datatypes.NewTime(18, 0, 0, 0),
			EndTime:   datatypes.NewTime(23, 0, 0, 0),
		},
	}

	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 1).Return(model.UserAvailability{UserID: 1, Availability: availability1, MeetingDurationMins: 30}, nil)
	suite.mockUserAvailabilityRepository.On("Get", suite.ctx, 2).Return(model.UserAvailability{}, errors.New("some error"))

	resp, err := suite.service.GetAvailabilityOverlap(suite.ctx, 1, 2)
	suite.Equal("some error", err.Error())
	suite.Equal(0, len(resp.Overlap))
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
