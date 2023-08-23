package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type UserTestSuite struct {
	suite.Suite
	controller  User
	mockService *MockUserService
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockService = &MockUserService{}
	suite.controller = User{
		userService: suite.mockService,
	}

}

func (suite *UserTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test","email":"test@example.xyz"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.User{Name: "test", Email: "test@example.xyz"}).Return(model.User{}, nil)

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(w.Result().StatusCode, http.StatusCreated)
	suite.Empty(body)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestCreateShouldReturnBadRequestWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(w.Result().StatusCode, http.StatusBadRequest)
	suite.Equal(`{"status_text":"bad request","message":"email is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "Created")
}

func (suite *UserTestSuite) TestCreateShouldReturnServerErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test","email":"test@example.xyz"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.User{Name: "test", Email: "test@example.xyz"}).Return(model.User{}, errors.New("some error"))

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(w.Result().StatusCode, http.StatusInternalServerError)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestSetAvailabilityHappyPath() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/1/availability", strings.NewReader(
		`{
			"availability":[
				{
					"day": "monday",
					"start_time": "10:00",
					"end_time": "17:00"
				},{
					"day": "tuesday",
					"start_time": "09:00",
					"end_time": "17:00"
				}
			],
			"meeting_duration_mins": 30
		}`))
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")

	suite.mockService.On("SetAvailability", req.Context(), 1, contract.UserAvailability{
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
	}).Return(model.UserAvailability{}, nil)

	suite.controller.SetAvailability(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(w.Result().StatusCode, http.StatusOK)
	suite.Empty(body)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestSetAvailabilityShouldReturnBadRequestWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPost, "/users/1/availability", strings.NewReader(`{"meeting_duration_mins": 30}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.controller.SetAvailability(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(w.Result().StatusCode, http.StatusBadRequest)
	suite.Equal(`{"status_text":"bad request","message":"at least one day's availability is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "SetAvailability")
}

func (suite *UserTestSuite) TestSetAvailabilityShouldReturnServerErrorWhenServiceReturnsError() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/1/availability", strings.NewReader(
		`{
			"availability":[
				{
					"day": "monday",
					"start_time": "10:00",
					"end_time": "17:00"
				},{
					"day": "tuesday",
					"start_time": "09:00",
					"end_time": "17:00"
				}
			],
			"meeting_duration_mins": 30
		}`))
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")

	suite.mockService.On("SetAvailability", req.Context(), 1, contract.UserAvailability{
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
	}).Return(model.UserAvailability{}, errors.New("some error"))

	suite.controller.SetAvailability(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(w.Result().StatusCode, http.StatusInternalServerError)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func TestUserTest(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
