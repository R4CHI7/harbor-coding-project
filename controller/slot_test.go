package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SlotTestSuite struct {
	suite.Suite
	controller      Slot
	mockSlotService *MockSlotService
}

func (suite *SlotTestSuite) SetupTest() {
	suite.mockSlotService = &MockSlotService{}
	suite.controller = NewSlot(suite.mockSlotService)
}

func (suite *SlotTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users/1/slot?num_days=14", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	w := httptest.NewRecorder()
	suite.mockSlotService.On("Create", req.Context(), 1, 14).Return(60, nil)

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusCreated, res.StatusCode)
	suite.Equal(`{"num_slots":60}
`, string(body))
	suite.mockSlotService.AssertExpectations(suite.T())
}

func (suite *SlotTestSuite) TestCreatShouldReturnErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users/1/slot?num_days=14", nil)
	req = req.WithContext(context.WithValue(context.Background(), ContextUserIDKey, 1))
	w := httptest.NewRecorder()
	suite.mockSlotService.On("Create", req.Context(), 1, 14).Return(-1, errors.New("some error"))

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
	suite.mockSlotService.AssertExpectations(suite.T())
}

func TestSlotTestSuite(t *testing.T) {
	suite.Run(t, new(SlotTestSuite))
}
