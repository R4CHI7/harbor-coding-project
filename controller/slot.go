package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/harbor-xyz/coding-project/contract"
)

type Slot struct {
	slotService SlotService
}

func (slot Slot) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(ContextUserIDKey).(int)

	numDaysParam := r.URL.Query().Get("num_days")
	if numDaysParam == "" {
		render.Render(w, r, contract.ErrorRenderer(errors.New("num_days is required")))
		return
	}
	numDays, err := strconv.Atoi(numDaysParam)
	if err != nil {
		render.Render(w, r, contract.ErrorRenderer(errors.New("invalid num_days")))
		return
	}

	numSlots, err := slot.slotService.Create(ctx, userID, numDays)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]int{"num_slots": numSlots})
}

func NewSlot(slotService SlotService) Slot {
	return Slot{slotService: slotService}
}
