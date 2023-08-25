package model

import (
	"time"

	"gorm.io/datatypes"
)

type Availability struct {
	StartTime datatypes.Time
	EndTime   datatypes.Time
}

type DayAvailability struct {
	Day       Day            `json:"day"`
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

func (availability UserAvailability) GetAvailabilityMap() map[Day]Availability {
	m := make(map[Day]Availability)
	for _, a := range availability.Availability {
		m[a.Day] = Availability{
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
		}
	}
	return m
}
