package event

import (
	"poc-event-source/internal/domain"
	"poc-event-source/internal/infrastructure/model"
)

func (e *eventRepository) CreateEvent(event *domain.EventSource) (*domain.EventSource, error) {
	ev := &model.EventSource{
		AggregateID: event.AggregateID,
		Type:        event.Type,
		Payload:     event.Payload,
	}
	if err := e.db.Create(ev).Error; err != nil {
		return nil, err
	}
	event.ID = ev.ID
	event.Version = ev.Version
	event.CreatedAt = ev.CreatedAt
	return event, nil
}
