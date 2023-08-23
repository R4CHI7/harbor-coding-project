package repository

import (
	"context"
	"log"

	"github.com/harbor-xyz/coding-project/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserAvailability struct {
	db *gorm.DB
}

func (availability UserAvailability) Set(ctx context.Context, input model.UserAvailability) (model.UserAvailability, error) {
	err := availability.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"availability": input.Availability, "meeting_duration_mins": input.MeetingDurationMins}),
	}).Create(&input).Error
	if err != nil {
		log.Printf("error occurred while saving user availability in DB: %s", err.Error())
		return model.UserAvailability{}, err
	}

	return input, nil
}

func (availability UserAvailability) Get(ctx context.Context, userID int) (model.UserAvailability, error) {
	ua := model.UserAvailability{}
	err := availability.db.Find(&ua, userID).Error
	if err != nil {
		log.Printf("error occurred while getting user availability from DB: %s", err.Error())
		return model.UserAvailability{}, err
	}

	return ua, nil
}

func NewUserAvailability(db *gorm.DB) UserAvailability {
	return UserAvailability{db: db}
}
