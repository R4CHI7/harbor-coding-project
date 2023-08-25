package service

import (
	"context"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/model"
)

type Event struct {
	eventRepository EventRepository
	slotRepository  SlotRepository
}

func (event Event) Create(ctx context.Context, userID int, input contract.Event) (contract.EventResponse, error) {
	slot, err := event.slotRepository.GetByID(ctx, input.SlotID)
	if err != nil {
		return contract.EventResponse{}, err
	}
	eventObj := model.Event{
		UserID:       uint(userID),
		SlotID:       uint(input.SlotID),
		InviteeEmail: input.InviteeEmail,
		InviteeName:  input.InviteeName,
		InviteeNotes: input.InviteeNotes,
		StartTime:    slot.StartTime,
		EndTime:      slot.EndTime,
	}

	eventObj, err = event.eventRepository.Create(ctx, eventObj)
	if err != nil {
		return contract.EventResponse{}, err
	}

	err = event.slotRepository.BookSlot(ctx, input.SlotID)
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
		StartTime:    eventObj.StartTime,
		EndTime:      eventObj.EndTime,
	}, nil
}

func (event Event) GetAll(ctx context.Context, userID int) (contract.EventListResponse, error) {
	events, err := event.eventRepository.GetAll(ctx, userID)
	if err != nil {
		return contract.EventListResponse{}, err
	}

	resp := make([]contract.EventResponse, 0)
	for _, eventObj := range events {
		resp = append(resp, contract.EventResponse{
			ID:           int(eventObj.ID),
			UserID:       int(eventObj.UserID),
			SlotID:       int(eventObj.SlotID),
			InviteeEmail: eventObj.InviteeEmail,
			InviteeName:  eventObj.InviteeName,
			InviteeNotes: eventObj.InviteeNotes,
			CreatedAt:    eventObj.CreatedAt,
			StartTime:    eventObj.StartTime,
			EndTime:      eventObj.EndTime,
		})
	}

	return contract.EventListResponse{Events: resp}, nil
}

func NewEvent(eventRepository EventRepository, slotRepository SlotRepository) Event {
	return Event{eventRepository: eventRepository, slotRepository: slotRepository}
}
