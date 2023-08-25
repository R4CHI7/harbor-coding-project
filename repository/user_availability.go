package repository

import (
	"context"
	"database/sql"
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
	res := availability.db.Find(&ua, userID)
	if res.Error != nil {
		log.Printf("error occurred while getting user availability from DB: %s", res.Error.Error())
		return model.UserAvailability{}, res.Error
	}

	if res.RowsAffected == 0 {
		log.Printf("user availability not found for user: %d", userID)
		return model.UserAvailability{}, sql.ErrNoRows
	}

	return ua, nil
}

func NewUserAvailability(db *gorm.DB) UserAvailability {
	return UserAvailability{db: db}
}
