package domain

import (
	"gorm.io/datatypes"
	"time"
)

type EventSourceEnum string

const (
	CreateUser EventSourceEnum = "CREATE_USER"
	DeleteUser EventSourceEnum = "DELETE_USER"
)

type EventSource struct {
	ID          string
	AggregateID string
	Type        string
	Payload     datatypes.JSON
	Version     int
	CreatedAt   time.Time
}
