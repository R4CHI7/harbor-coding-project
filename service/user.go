package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
)

type User struct {
	userRepository         UserRepository
	availabilityRepository UserAvailabilityRepository
}

func (user User) Create(ctx context.Context, input contract.User) (model.User, error) {
	userObj := model.User{
		Name:  input.Name,
		Email: input.Email,
	}

	return user.userRepository.Create(ctx, userObj)
}

func (user User) SetAvailability(ctx context.Context, userID int, input contract.UserAvailability) (model.UserAvailability, error) {
	availabilityObj := model.UserAvailability{
		UserID:              uint(userID),
		Availability:        input.Availability,
		MeetingDurationMins: input.MeetingDurationMins,
	}

	return user.availabilityRepository.Set(ctx, availabilityObj)
}

func NewUser(userRepository UserRepository, availabilityRepository UserAvailabilityRepository) User {
	return User{userRepository: userRepository, availabilityRepository: availabilityRepository}
}
