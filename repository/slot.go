package repository

import (
	"context"
	"log"
	"time"

	"github.com/harbor-xyz/coding-project/model"
	"gorm.io/gorm"
)

type Slot struct {
	db *gorm.DB
}

func (slot Slot) Get(ctx context.Context, userID int, startTimeThreshold, endTimeThreshold time.Time) ([]model.Slot, error) {
	slots := make([]model.Slot, 0)
	err := slot.db.Find(&slots, "user_id = $1 AND start_time BETWEEN $2 AND $3", userID, startTimeThreshold, endTimeThreshold).Error
	if err != nil {
		log.Printf("error occurred while fetching slots for user: %s", err.Error())
		return nil, err
	}
	return slots, nil
}

func (slot Slot) Create(ctx context.Context, slots []model.Slot) error {
	err := slot.db.Create(slots).Error
	if err != nil {
		log.Printf("error occurred while inserting slots: %s", err.Error())
	}
	return err
}

func NewSlot(db *gorm.DB) Slot {
	return Slot{db: db}
}
