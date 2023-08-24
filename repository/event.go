package repository

import (
	"context"
	"log"

	"github.com/harbor-xyz/coding-project/model"
	"gorm.io/gorm"
)

type Event struct {
	db *gorm.DB
}

func (event Event) Create(ctx context.Context, obj model.Event) (model.Event, error) {
	err := event.db.Create(&obj).Error
	if err != nil {
		log.Printf("error occurred while saving event in DB: %s", err.Error())
		return model.Event{}, err
	}

	return obj, nil
}

func (event Event) Get(ctx context.Context, userID int) ([]model.Event, error) {
	events := make([]model.Event, 0)
	err := event.db.Find(&events, "user_id = $1", userID).Error
	if err != nil {
		log.Printf("error occurred while saving event in DB: %s", err.Error())
		return nil, err
	}

	return events, nil
}

func NewEvent(db *gorm.DB) Event {
	return Event{db: db}
}
