package model

import (
	"errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type EventSource struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	AggregateID string         `gorm:"type:uuid;index"`
	Type        string         `gorm:"index;not null;default:nul"`
	Payload     datatypes.JSON `gorm:"type:jsonb;not null"`
	Version     int            `gorm:"not null;default:1"`
	CreatedAt   time.Time      `gorm:"not null"`
}

func (e *EventSource) BeforeSave(_ *gorm.DB) (err error) {
	if e.Type == "" {
		return errors.New(" null value in column \"type\" of relation \"event_sources\" violates not-null constraint")
	}

	return nil
}
