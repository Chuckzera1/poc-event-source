package event

import "poc-event-source/internal/infrastructure/model"

func (e *eventRepository) CreateEvent(event *model.EventSource) (*model.EventSource, error) {
	ev := event
	err := e.db.Create(ev).Error
	if err != nil {
		return nil, err
	}

	return ev, err
}
