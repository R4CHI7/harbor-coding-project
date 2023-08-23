package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
)

type User struct {
	userRepository UserRepository
}

func (user User) Create(ctx context.Context, input contract.User) (model.User, error) {
	userObj := model.User{
		Name:  input.Name,
		Email: input.Email,
	}

	insertedObj, err := user.userRepository.Create(ctx, userObj)
	if err != nil {
		return model.User{}, err
	}

	return insertedObj, nil
}

func NewUser(userRepository UserRepository) User {
	return User{userRepository: userRepository}
}
