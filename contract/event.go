package contract

import (
	"errors"
	"net/http"
	"time"
)

type Event struct {
	SlotID       int    `json:"slot_id"`
	InviteeEmail string `json:"invitee_email"`
	InviteeName  string `json:"invitee_name"`
	InviteeNotes string `json:"invitee_notes"`
}

func (event *Event) Bind(r *http.Request) error {
	if event.SlotID == 0 {
		return errors.New("slot_id is required")
	}

	if event.InviteeEmail == "" {
		return errors.New("invitee_email is required")
	}

	if event.InviteeName == "" {
		return errors.New("invitee_name is required")
	}

	return nil
}

type EventResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	SlotID       int       `json:"slot_id"`
	InviteeEmail string    `json:"invitee_email"`
	InviteeName  string    `json:"invitee_name"`
	InviteeNotes string    `json:"invitee_notes"`
	CreatedAt    time.Time `json:"created_at"`
}
