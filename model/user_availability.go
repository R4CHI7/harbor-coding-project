package model

import (
	"time"

	"gorm.io/datatypes"
)

type DayAvailability struct {
	Day       string         `json:"day"`
	StartTime datatypes.Time `json:"start_time"`
	EndTime   datatypes.Time `json:"end_time"`
}

type UserAvailability struct {
	UserID              uint `gorm:"uniqueIndex"`
	Availability        datatypes.JSONSlice[DayAvailability]
	MeetingDurationMins int
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}
