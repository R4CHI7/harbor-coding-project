package service

import (
	"context"
	"time"

	"github.com/harbor-xyz/coding-project/model"
)

type UserRepository interface {
	Create(context.Context, model.User) (model.User, error)
}

type UserAvailabilityRepository interface {
	Set(context.Context, model.UserAvailability) (model.UserAvailability, error)
	Get(context.Context, int) (model.UserAvailability, error)
}

type SlotRepository interface {
	Create(context.Context, []model.Slot) error
	Get(context.Context, int, time.Time, time.Time) ([]model.Slot, error)
	GetByID(context.Context, int) (model.Slot, error)
	DeleteByID(context.Context, int) error
	BookSlot(context.Context, int) error
}

type EventRepository interface {
	Create(context.Context, model.Event) (model.Event, error)
	GetAll(context.Context, int) ([]model.Event, error)
}
