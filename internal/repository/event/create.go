package event

import (
	"context"

	"poc-event-source/internal/domain"
	"poc-event-source/internal/infrastructure/model"
)

func (e *eventRepository) CreateEvent(ctx context.Context, event *domain.EventSource) (*domain.EventSource, error) {
	ev := &model.EventSource{
		AggregateID: event.AggregateID,
		Type:        event.Type,
		Payload:     event.Payload,
	}
	if err := e.db.WithContext(ctx).Create(ev).Error; err != nil {
		return nil, err
	}
	event.ID = ev.ID
	event.Version = ev.Version
	event.CreatedAt = ev.CreatedAt
	return event, nil
}
