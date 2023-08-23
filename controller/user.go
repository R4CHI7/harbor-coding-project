package controller

import (
	"log"
	"net/http"

	"github.com/go-chi/render"

	"github.com/harbor-xyz/coding-project/contract"
)

type User struct {
	userService UserService
}

func (user User) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.User{}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	_, err := user.userService.Create(ctx, input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func NewUser(userService UserService) User {
	return User{
		userService: userService,
	}
}
