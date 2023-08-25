package controller

import (
	"log"
	"net/http"

	"github.com/go-chi/render"

	"github.com/harbor-xyz/coding-project/contract"
)

type Event struct {
	eventService EventService
}

// Create - Creates a new event
// @Summary This API creates a new event for the user with invitee details.
// @Tags event
// @Accept  json
// @Produce  json
// @Param event body contract.Event true "Add event"
// @Param user_id path int true "user id"
// @Success 200 {object} contract.EventResponse
// @Router /users/{user_id}/events [post]
func (event Event) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.Event{}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	userID := ctx.Value(ContextUserIDKey).(int)

	resp, err := event.eventService.Create(ctx, userID, input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

// GetAll - Returns events for user
// @Summary This API returns all events for a given user ID.
// @Tags event
// @Accept  json
// @Produce  json
// @Param user_id path int true "user id"
// @Success 200 {object} contract.EventListResponse
// @Router /users/{user_id}/events [get]
func (event Event) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(ContextUserIDKey).(int)

	resp, err := event.eventService.GetAll(ctx, userID)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.JSON(w, r, resp)
}

func NewEvent(eventService EventService) Event {
	return Event{eventService: eventService}
}
