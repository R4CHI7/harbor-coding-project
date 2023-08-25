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
	err := slot.db.Order("id").Find(&slots, "user_id = $1 AND start_time BETWEEN $2 AND $3", userID, startTimeThreshold, endTimeThreshold).Error
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

func (slot Slot) GetByID(ctx context.Context, slotID int) (model.Slot, error) {
	slotObj := model.Slot{ID: uint(slotID)}
	err := slot.db.Find(&slotObj).Error
	if err != nil {
		log.Printf("error occurred while fetching slot from db: %s", err.Error())
		return model.Slot{}, err
	}
	return slotObj, nil
}

func (slot Slot) DeleteByID(ctx context.Context, slotID int) error {
	slotObj := model.Slot{ID: uint(slotID)}
	err := slot.db.Model(&slotObj).Updates(model.Slot{Status: model.StatusDeleted, DeletedAt: time.Now()}).Error
	if err != nil {
		log.Printf("error occurred while deleting slot from db: %s", err.Error())
		return err
	}
	return nil
}

func (slot Slot) BookSlot(ctx context.Context, slotID int) error {
	slotObj := model.Slot{ID: uint(slotID)}
	err := slot.db.Model(&slotObj).Update("status", model.StatusBooked).Error
	if err != nil {
		log.Printf("error occurred while booking slot in db: %s", err.Error())
		return err
	}
	return nil
}

func NewSlot(db *gorm.DB) Slot {
	return Slot{db: db}
}
