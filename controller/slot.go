package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/harbor-xyz/coding-project/contract"
)

type Slot struct {
	slotService SlotService
}

// Create - Creates slots for a user
// @Summary This API creates slots for a user for given number of days.
// @Tags slot
// @Accept  json
// @Produce  json
// @Param num_days query int true "number of days to create slots"
// @Param user_id path int true "user id"
// @Router /users/{user_id}/slots [post]
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

// GetAll - Gets slots for a user
// @Summary This API returns slots for a user starting today till 14 days.
// @Tags slot
// @Accept  json
// @Produce  json
// @Param user_id path int true "user id"
// @Success 200 {object} contract.SlotList
// @Router /users/{user_id}/slots [get]
func (slot Slot) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(ContextUserIDKey).(int)

	slots, err := slot.slotService.GetAll(ctx, userID)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.JSON(w, r, slots)
}

// Delete - Deletes slot by ID
// @Summary This API returns deletes a slot by ID.
// @Tags slot
// @Accept  json
// @Produce  json
// @Param user_id path int true "user id"
// @Param slot_id path int true "slot id"
// @Router /users/{user_id}/slots/{slot_id} [delete]
func (slot Slot) Delete(w http.ResponseWriter, r *http.Request) {
	slotID := chi.URLParam(r, "slotID")
	if slotID == "" {
		render.Render(w, r, contract.ErrorRenderer(errors.New("user ID is required")))
		return
	}
	id, err := strconv.Atoi(slotID)
	if err != nil {
		render.Render(w, r, contract.ErrorRenderer(errors.New("invalid user ID")))
	}

	err = slot.slotService.DeleteByID(r.Context(), id)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
}

func NewSlot(slotService SlotService) Slot {
	return Slot{slotService: slotService}
}
