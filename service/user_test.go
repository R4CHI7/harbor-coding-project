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

type UserTestSuite struct {
	suite.Suite
	service            User
	mockUserRepository *MockUserRepository
	ctx                context.Context
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockUserRepository = &MockUserRepository{}
	suite.service = User{
		userRepository: suite.mockUserRepository,
	}
	suite.ctx = context.Background()
}

func (suite *UserTestSuite) TestHappyFlow() {
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

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
