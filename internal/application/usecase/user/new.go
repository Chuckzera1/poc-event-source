package user

import (
	"context"

	"poc-event-source/internal/application/dto"
	eventhandler "poc-event-source/internal/application/usecase/event"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, input dto.CreateUserReqDTO) error
}

type createUserUseCase struct {
	eventHandler eventhandler.MainHandlerUseCase
}

func NewCreateUserUseCase(eventHandler eventhandler.MainHandlerUseCase) CreateUserUseCase {
	return &createUserUseCase{eventHandler: eventHandler}
}
