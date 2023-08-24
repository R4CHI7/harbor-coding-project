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

func (user User) Create(ctx context.Context, input contract.User) (contract.UserResponse, error) {
	userObj := model.User{
		Name:  input.Name,
		Email: input.Email,
	}

	userObj, err := user.userRepository.Create(ctx, userObj)
	if err != nil {
		return contract.UserResponse{}, err
	}

	return contract.UserResponse{ID: userObj.ID}, nil
}

func (user User) SetAvailability(ctx context.Context, userID int, input contract.UserAvailability) (model.UserAvailability, error) {
	availabilityObj := model.UserAvailability{
		UserID:              uint(userID),
		Availability:        input.Availability,
		MeetingDurationMins: input.MeetingDurationMins,
	}

	return user.availabilityRepository.Set(ctx, availabilityObj)
}

func (user User) GetAvailability(ctx context.Context, userID int) (contract.UserAvailability, error) {
	availability, err := user.availabilityRepository.Get(ctx, userID)
	if err != nil {
		return contract.UserAvailability{}, err
	}

	return contract.UserAvailability{
		Availability:        availability.Availability,
		MeetingDurationMins: availability.MeetingDurationMins,
	}, nil
}

func NewUser(userRepository UserRepository, availabilityRepository UserAvailabilityRepository) User {
	return User{userRepository: userRepository, availabilityRepository: availabilityRepository}
}
