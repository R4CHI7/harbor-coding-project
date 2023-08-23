package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/model"
)

type UserRepository interface {
	Create(context.Context, model.User) (model.User, error)
}
