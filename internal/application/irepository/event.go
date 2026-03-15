package irepository

import (
	"context"

	"poc-event-source/internal/domain"
)

type CreateEventRepository interface {
	CreateEvent(ctx context.Context, event *domain.EventSource) (*domain.EventSource, error)
}

type EventRepository interface {
	CreateEventRepository
}
