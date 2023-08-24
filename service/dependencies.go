package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/model"
)

type UserRepository interface {
	Create(context.Context, model.User) (model.User, error)
}

type UserAvailabilityRepository interface {
	Set(context.Context, model.UserAvailability) (model.UserAvailability, error)
	Get(context.Context, int) (model.UserAvailability, error)
}
