package model

import (
	"time"
)

type status int

const (
	StatusCreated status = 0
	StatusBooked  status = 1
	StatusDeleted status = 2
	StatusExpired status = 4
)

func (s status) String() string {
	switch s {
	case StatusCreated:
		return "created"
	case StatusBooked:
		return "booked"
	case StatusDeleted:
		return "deleted"
	case StatusExpired:
		return "expired"
	}
	return ""
}

type Slot struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	StartTime time.Time
	EndTime   time.Time
	Status    status
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time

	Event Event
}
