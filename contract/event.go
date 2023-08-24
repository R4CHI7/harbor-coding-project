package contract

import "time"

type Event struct {
	SlotID       int    `json:"slot_id"`
	InviteeEmail string `json:"invitee_email"`
	InviteeName  string `json:"invitee_name"`
	InviteeNotes string `json:"invitee_notes"`
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
