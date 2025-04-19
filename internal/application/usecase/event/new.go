package event

import (
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/application/irepository"
)

type MainHandlerUseCase interface {
	Handler(event dto.EventReqDTO) error
}
type mainHandlerUseCase struct {
	createEventRepo irepository.CreateEventRepository
}

func NewMainHandler(createEventRepo irepository.CreateEventRepository) MainHandlerUseCase {
	return &mainHandlerUseCase{
		createEventRepo: createEventRepo,
	}
}
