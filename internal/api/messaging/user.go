package messaging

import (
	"context"
	"encoding/json"
	"log"

	"poc-event-source/internal/application"
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/infrastructure/model"
)

func (ub *UserBroker) Subscribe() error {
	ctx := context.Background()
	_, err := ub.broker.Subscribe(ctx, "user", func(ctx context.Context, msg *application.Message) {
		var envelope dto.EventMessage
		if err := json.Unmarshal(msg.Data, &envelope); err != nil {
			log.Printf("Subscribe: error unmarshaling message: %v", err)
			msg.Ack()
			return
		}

		handler, ok := ub.handlers[envelope.Type]
		if !ok {
			log.Printf("Subscribe: no handler for event type: %s", envelope.Type)
			msg.Ack()
			return
		}

		handler(ctx, msg)
	})
	return err
}

func (ub *UserBroker) handleCreate(ctx context.Context, msg *application.Message) {
	var envelope dto.EventMessage
	if err := json.Unmarshal(msg.Data, &envelope); err != nil {
		log.Printf("handleCreate: error unmarshaling envelope: %v", err)
		msg.Ack()
		return
	}

	var input dto.CreateUserReqDTO
	if err := json.Unmarshal(envelope.Payload, &input); err != nil {
		log.Printf("handleCreate: error unmarshaling payload: %v", err)
		msg.Ack()
		return
	}

	hashed, err := ub.pwdUtil.HashPassword(input.Password)
	if err != nil {
		log.Printf("handleCreate: error hashing password: %v", err)
		msg.Ack()
		return
	}

	if _, err = ub.userRepo.CreateUser(&model.User{Username: input.Username, Password: hashed}); err != nil {
		log.Printf("handleCreate: error creating user: %v", err)
		return
	}

	msg.Ack()
}
