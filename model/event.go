package model

import "time"

type Event struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	SlotID       uint
	InviteeEmail string `gorm:"not null"`
	InviteeName  string `gorm:"not null"`
	InviteeNotes string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	DeletedAt    time.Time
}
