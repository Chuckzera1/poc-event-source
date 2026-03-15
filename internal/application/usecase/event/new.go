package event

import (
	"context"

	"poc-event-source/internal/application"
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/application/irepository"
)

type MainHandlerUseCase interface {
	Handler(ctx context.Context, topic string, event dto.EventReqDTO) error
}

type mainHandlerUseCase struct {
	createEventRepo irepository.CreateEventRepository
	broker          application.Broker
}

func NewMainHandler(createEventRepo irepository.CreateEventRepository, broker application.Broker) MainHandlerUseCase {
	return &mainHandlerUseCase{
		createEventRepo: createEventRepo,
		broker:          broker,
	}
}
