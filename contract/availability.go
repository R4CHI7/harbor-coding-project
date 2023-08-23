package contract

import (
	"errors"
	"net/http"

	"github.com/harbor-xyz/coding-project/model"
)

type UserAvailability struct {
	Availability        []model.DayAvailability `json:"availability"`
	MeetingDurationMins int                     `json:"meeting_duration_mins"`
}

func (availability *UserAvailability) Bind(r *http.Request) error {
	if len(availability.Availability) == 0 {
		return errors.New("at least one day's availability is required")
	}

	if availability.MeetingDurationMins < 15 {
		return errors.New("meeting_duration should be at least 15")
	}

	return nil
}
