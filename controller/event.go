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

func (event Event) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(ContextUserIDKey).(int)

	resp, err := event.eventService.Get(ctx, userID)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.JSON(w, r, resp)
}

func NewEvent(eventService EventService) Event {
	return Event{eventService: eventService}
}
