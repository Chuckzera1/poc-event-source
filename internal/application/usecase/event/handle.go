package event

import (
	"context"
	"encoding/json"

	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/domain"
)

func (m *mainHandlerUseCase) Handler(ctx context.Context, topic string, event dto.EventReqDTO) error {
	_, err := m.createEventRepo.CreateEvent(ctx, &domain.EventSource{
		Type:    event.Type,
		Payload: event.Payload,
	})
	if err != nil {
		return err
	}

	envelope, err := json.Marshal(dto.EventMessage{
		Type:    event.Type,
		Payload: json.RawMessage(event.Payload),
	})
	if err != nil {
		return err
	}

	return m.broker.Publish(ctx, topic, envelope)
}
