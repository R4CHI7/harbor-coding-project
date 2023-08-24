package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

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

	resp, err := user.userService.Create(ctx, input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func (user User) SetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.UserAvailability{}
	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	userID := ctx.Value(ContextUserIDKey).(int)

	_, err := user.userService.SetAvailability(ctx, userID, input)

	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (user User) GetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(ContextUserIDKey).(int)

	availability, err := user.userService.GetAvailability(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, contract.NotFoundErrorRenderer(err))
			return
		}
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.JSON(w, r, availability)
}

func (user User) GetAvailabilityOverlap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user1ID := ctx.Value(ContextUserIDKey).(int)
	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		render.Render(w, r, contract.ErrorRenderer(errors.New("user ID is required")))
		return
	}
	user2ID, err := strconv.Atoi(userIDParam)
	if err != nil {
		render.Render(w, r, contract.ErrorRenderer(errors.New("invalid user ID")))
		return
	}

	if user1ID == user2ID {
		render.Render(w, r, contract.ErrorRenderer(errors.New("user IDs should be different")))
		return
	}

	overlap, err := user.userService.GetAvailabilityOverlap(ctx, user1ID, user2ID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, contract.NotFoundErrorRenderer(err))
			return
		}
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.JSON(w, r, overlap)

}

func NewUser(userService UserService) User {
	return User{
		userService: userService,
	}
}
