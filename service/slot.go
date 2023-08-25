package service

import (
	"context"
	"errors"
	"time"

	"github.com/harbor-xyz/coding-project/model"
)

type Slot struct {
	slotRepository         SlotRepository
	availabilityRepository UserAvailabilityRepository
}

func (slot Slot) Create(ctx context.Context, userID, numDays int) (int, error) {
	now := time.Now()
	// Check if slots already exist in the given time period
	slots, err := slot.slotRepository.Get(ctx, userID, now, now.AddDate(0, 0, numDays))
	if err != nil {
		return -1, err
	}

	if len(slots) > 0 {
		return -1, errors.New("slots already exist")
	}

	// Get availability for the user
	availability, err := slot.availabilityRepository.Get(ctx, userID)
	if err != nil {
		return -1, err
	}
	availabilityMap := availability.GetAvailabilityMap()
	meetingDuration := availability.MeetingDurationMins

	// Prepare slots based on days and hours of availability and meeting duration
	slots = make([]model.Slot, 0)
	for i := 0; i < numDays; i++ {
		t := now.AddDate(0, 0, i)
		day := model.GetDayFromInt(int(t.Weekday()))
		// If availability exists for this day
		if availability, exists := availabilityMap[day]; exists {
			s, _ := time.Parse("15:04:05", availability.StartTime.String())
			e, _ := time.Parse("15:04:05", availability.EndTime.String())

			startTime := time.Date(t.Year(), t.Month(), t.Day(), s.Hour(), s.Minute(), s.Second(), s.Nanosecond(), t.Location())
			endTime := time.Date(t.Year(), t.Month(), t.Day(), e.Hour(), e.Minute(), e.Second(), e.Nanosecond(), t.Location())
			for startTime.Before(endTime) {
				end := startTime.Add(time.Minute * time.Duration(meetingDuration))
				slots = append(slots, model.Slot{
					UserID:    uint(userID),
					StartTime: startTime,
					EndTime:   end,
					Status:    model.StatusCreated,
				})
				startTime = end
			}
		}
	}
	// Insert slots
	err = slot.slotRepository.Create(ctx, slots)
	if err != nil {
		return -1, nil
	}
	return len(slots), nil
}

func NewSlot(slotRepository SlotRepository, availabilityRepository UserAvailabilityRepository) Slot {
	return Slot{slotRepository: slotRepository, availabilityRepository: availabilityRepository}
}
