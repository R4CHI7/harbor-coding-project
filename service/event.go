package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
)

type Event struct {
	eventRepository EventRepository
}

func (event Event) Create(ctx context.Context, userID int, input contract.Event) (contract.EventResponse, error) {
	eventObj := model.Event{
		UserID:       uint(userID),
		SlotID:       uint(input.SlotID),
		InviteeEmail: input.InviteeEmail,
		InviteeName:  input.InviteeName,
		InviteeNotes: input.InviteeNotes,
	}

	eventObj, err := event.eventRepository.Create(ctx, eventObj)
	if err != nil {
		return contract.EventResponse{}, err
	}

	return contract.EventResponse{
		ID:           int(eventObj.ID),
		UserID:       int(eventObj.UserID),
		SlotID:       int(eventObj.SlotID),
		InviteeEmail: eventObj.InviteeEmail,
		InviteeName:  eventObj.InviteeName,
		InviteeNotes: eventObj.InviteeNotes,
		CreatedAt:    eventObj.CreatedAt,
	}, nil
}

func NewEvent(eventRepository EventRepository) Event {
	return Event{eventRepository: eventRepository}
}
