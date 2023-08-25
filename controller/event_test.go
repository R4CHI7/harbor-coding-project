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
	req := httptest.NewRequest(http.MethodPost, "/users/1/event",
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

func (suite *EventTestSuite) TestCreatShouldReturnErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users/1/event",
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

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
