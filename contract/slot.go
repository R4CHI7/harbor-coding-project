package contract

import "time"

type Slot struct {
	ID        int       `json:"id"`
	UserID    uint      `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type SlotList struct {
	Slots []Slot `json:"slots"`
}
