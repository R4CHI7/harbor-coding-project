package controller

import (
	"context"
	"database/sql"
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
	suite.controller = NewUser(suite.mockService)

}

func (suite *UserTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test","email":"test@example.xyz"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.User{Name: "test", Email: "test@example.xyz"}).Return(contract.UserResponse{ID: 1}, nil)

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusCreated, res.StatusCode)
	suite.Equal(`{"id":1}
`, string(body))
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
	suite.Equal(http.StatusBadRequest, res.StatusCode)
	suite.Equal(`{"status_text":"bad request","message":"email is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "Created")
}

func (suite *UserTestSuite) TestCreateShouldReturnServerErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test","email":"test@example.xyz"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.User{Name: "test", Email: "test@example.xyz"}).Return(contract.UserResponse{}, errors.New("some error"))

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
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
	suite.Equal(http.StatusOK, res.StatusCode)
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
	suite.Equal(http.StatusBadRequest, res.StatusCode)
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
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestGetAvailabilityHappyPath() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/availability", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockService.On("GetAvailability", req.Context(), 1).Return(contract.UserAvailability{
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

	suite.controller.GetAvailability(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusOK, w.Result().StatusCode)
	suite.Equal(`{"availability":[{"day":"monday","start_time":"10:00:00","end_time":"17:00:00"},{"day":"tuesday","start_time":"09:00:00","end_time":"17:00:00"}],"meeting_duration_mins":30}
`, string(body))
}

func (suite *UserTestSuite) TestGetAvailabilityReturnsServerErrorWhenServiceReturnsError() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/availability", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockService.On("GetAvailability", req.Context(), 1).Return(contract.UserAvailability{}, errors.New("some error"))

	suite.controller.GetAvailability(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusInternalServerError, w.Result().StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body))
}

func (suite *UserTestSuite) TestGetAvailabilityReturnsNotFoundErrorWhenAvailabilityDoesNotExist() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/availability", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockService.On("GetAvailability", req.Context(), 1).Return(contract.UserAvailability{}, sql.ErrNoRows)

	suite.controller.GetAvailability(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusNotFound, w.Result().StatusCode)
	suite.Equal(`{"status_text":"not found","message":"sql: no rows in result set"}
`, string(body))
}

func (suite *UserTestSuite) TestGetAvailabilityOverlapHappyPath() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/availability_overlap?second_user_id=2", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockService.On("GetAvailabilityOverlap", req.Context(), 1, 2).Return(contract.UserAvailabilityOverlap{
		Overlap: []model.DayAvailability{
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
	}, nil)

	suite.controller.GetAvailabilityOverlap(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusOK, w.Result().StatusCode)
	suite.Equal(`{"overlap":[{"day":"monday","start_time":"10:00:00","end_time":"17:00:00"},{"day":"tuesday","start_time":"09:00:00","end_time":"17:00:00"}]}
`, string(body))
}

func (suite *UserTestSuite) TestGetAvailabilityOverlapReturnsNullResponseWhenNoOverlapExists() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/availability_overlap?second_user_id=2", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockService.On("GetAvailabilityOverlap", req.Context(), 1, 2).Return(contract.UserAvailabilityOverlap{}, nil)

	suite.controller.GetAvailabilityOverlap(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusOK, w.Result().StatusCode)
	suite.Equal(`{"overlap":null}
`, string(body))
}

func (suite *UserTestSuite) TestGetAvailabilityOverlapReturnsErrorWhenServiceReturnsError() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/availability_overlap?second_user_id=2", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockService.On("GetAvailabilityOverlap", req.Context(), 1, 2).Return(contract.UserAvailabilityOverlap{}, errors.New("some error"))

	suite.controller.GetAvailabilityOverlap(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusInternalServerError, w.Result().StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body))
}

func TestUserTest(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
