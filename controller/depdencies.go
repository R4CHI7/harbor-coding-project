package controller

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
)

type UserService interface {
	Create(context.Context, contract.User) (model.User, error)
}
