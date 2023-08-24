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

func (user User) GetAvailabilityOverlap(ctx context.Context, userID1, userID2 int) (contract.UserAvailabilityOverlap, error) {
	availability1, err := user.availabilityRepository.Get(ctx, userID1)
	if err != nil {
		return contract.UserAvailabilityOverlap{}, err
	}
	availability2, err := user.availabilityRepository.Get(ctx, userID2)
	if err != nil {
		return contract.UserAvailabilityOverlap{}, err
	}

	overlap := make([]model.DayAvailability, 0)

	av1Map := availability1.GetAvailabilityMap()
	av2Map := availability2.GetAvailabilityMap()

	for day, availability := range av1Map {
		if a, exists := av2Map[day]; exists {
			// If there is an overlap between the 2 availabilities
			if availability.StartTime <= a.EndTime && availability.EndTime >= a.StartTime {
				overlappingAvailability := model.DayAvailability{Day: day}

				// Start time will be the greater of the 2 starts
				if availability.StartTime <= a.StartTime {
					overlappingAvailability.StartTime = a.StartTime
				} else {
					overlappingAvailability.StartTime = availability.StartTime
				}

				// End time will be the lower of the 2 ends
				if availability.EndTime <= a.EndTime {
					overlappingAvailability.EndTime = availability.EndTime
				} else {
					overlappingAvailability.EndTime = a.EndTime
				}

				overlap = append(overlap, overlappingAvailability)
			}
		}
	}

	return contract.UserAvailabilityOverlap{
		Overlap: overlap,
	}, nil
}

func NewUser(userRepository UserRepository, availabilityRepository UserAvailabilityRepository) User {
	return User{userRepository: userRepository, availabilityRepository: availabilityRepository}
}
