package controller

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
)

type UserService interface {
	Create(context.Context, contract.User) (contract.UserResponse, error)
	SetAvailability(context.Context, int, contract.UserAvailability) (model.UserAvailability, error)
	GetAvailability(context.Context, int) (contract.UserAvailability, error)
	GetAvailabilityOverlap(context.Context, int, int) (contract.UserAvailabilityOverlap, error)
}

type EventService interface {
	Create(context.Context, int, contract.Event) (contract.EventResponse, error)
}

type SlotService interface {
	Create(context.Context, int, int) (int, error)
}
