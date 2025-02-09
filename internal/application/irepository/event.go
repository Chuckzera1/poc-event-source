package irepository

import "poc-event-source/internal/infrastructure/model"

type CreateEventRepository interface {
	CreateEvent(event *model.EventSource) (*model.EventSource, error)
}

type EventRepository interface {
	CreateEventRepository
}
