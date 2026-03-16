package user

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/domain"
)

func (u *createUserUseCase) Execute(ctx context.Context, input dto.CreateUserReqDTO) error {
	payloadBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return u.eventHandler.Handler(ctx, "user", dto.EventReqDTO{
		Type:    string(domain.CreateUser),
		Payload: datatypes.JSON(payloadBytes),
	})
}
