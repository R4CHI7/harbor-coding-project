package model

import "time"

type Event struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	SlotID       uint
	InviteeEmail string `gorm:"not null"`
	InviteeName  string `gorm:"not null"`
	InviteeNotes string
	StartTime    time.Time `gorm:"not null"`
	EndTime      time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	DeletedAt    time.Time
}
