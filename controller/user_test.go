package controller

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/suite"
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

func (suite *UserTestSuite) TestHappyFlow() {
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

func (suite *UserTestSuite) TestCreateShouldReturnBadRequestWhenRequestBodyHasError() {
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

func TestUserTest(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
