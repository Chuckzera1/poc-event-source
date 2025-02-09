package event

import (
	"gorm.io/gorm"
	"poc-event-source/internal/application/irepository"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) irepository.EventRepository {
	return &eventRepository{db: db}
}
