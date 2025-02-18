package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/harbor-xyz/coding-project/contract"

	"github.com/stretchr/testify/suite"
)

type EventTestSuite struct {
	suite.Suite
	controller       Event
	mockEventService *MockEventService
}

func (suite *EventTestSuite) SetupTest() {
	suite.mockEventService = &MockEventService{}
	suite.controller = NewEvent(suite.mockEventService)
}

func (suite *EventTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users/1/events",
		strings.NewReader(`{"slot_id":1,"invitee_email":"test@example.xyz","invitee_name":"test","invitee_notes":"test"}`))
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockEventService.On("Create", req.Context(), 1, contract.Event{
		SlotID: 1, InviteeName: "test", InviteeEmail: "test@example.xyz", InviteeNotes: "test"}).
		Return(contract.EventResponse{ID: 1, UserID: 1, SlotID: 1, InviteeEmail: "test@example.xyz", InviteeName: "test", InviteeNotes: "test"}, nil)

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	suite.Equal(http.StatusCreated, res.StatusCode)
	suite.mockEventService.AssertExpectations(suite.T())
}

func (suite *EventTestSuite) TestCreateShouldReturnErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users/1/events",
		strings.NewReader(`{"slot_id":1,"invitee_email":"test@example.xyz","invitee_name":"test","invitee_notes":"test"}`))
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockEventService.On("Create", req.Context(), 1, contract.Event{
		SlotID: 1, InviteeName: "test", InviteeEmail: "test@example.xyz", InviteeNotes: "test"}).
		Return(contract.EventResponse{}, errors.New("some error"))

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
	suite.mockEventService.AssertExpectations(suite.T())
}

func (suite *EventTestSuite) TestGetAllHappyPath() {
	now := time.Now()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/events", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	suite.mockEventService.On("GetAll", req.Context(), 1).Return(contract.EventListResponse{
		Events: []contract.EventResponse{
			{
				ID:           1,
				UserID:       1,
				SlotID:       1,
				InviteeEmail: "test@example.xyz",
				InviteeName:  "test",
				StartTime:    now,
				EndTime:      now.Add(30 * time.Minute),
				CreatedAt:    now,
			},
		},
	}, nil)

	suite.controller.GetAll(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusOK, w.Result().StatusCode)
	suite.NotEmpty(body)
	suite.mockEventService.AssertExpectations(suite.T())
}

func (suite *EventTestSuite) TestGetAllReturnsServerErrorWhenServiceReturnsError() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1/events", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	req.Header.Add("Content-Type", "application/json")
	suite.mockEventService.On("GetAll", req.Context(), 1).Return(contract.EventListResponse{}, errors.New("some error"))

	suite.controller.GetAll(w, req)

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

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
