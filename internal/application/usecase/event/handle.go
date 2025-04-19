package event

import (
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/infrastructure/model"
)

func (m *mainHandlerUseCase) Handler(event dto.EventReqDTO) error {
	modelEv := &model.EventSource{
		Payload: event.Payload,
		Type:    event.Type,
	}

	_, err := m.createEventRepo.CreateEvent(modelEv)
	if err != nil {
		return err
	}

	return nil
}
