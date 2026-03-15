package irepository

import "poc-event-source/internal/domain"

type CreateEventRepository interface {
	CreateEvent(event *domain.EventSource) (*domain.EventSource, error)
}

type EventRepository interface {
	CreateEventRepository
}
