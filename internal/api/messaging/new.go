package messaging

import (
	"context"

	"poc-event-source/internal/api"
	"poc-event-source/internal/application"
	"poc-event-source/internal/application/irepository"
	"poc-event-source/internal/domain"
)

type UserBroker struct {
	broker   application.Broker
	userRepo irepository.UserRepository
	handlers map[string]func(context.Context, *application.Message)
}

func NewUserBroker(
	broker application.Broker,
	userRepo irepository.UserRepository,
) api.Subscriber {
	ub := &UserBroker{
		broker:   broker,
		userRepo: userRepo,
	}
	ub.handlers = map[string]func(context.Context, *application.Message){
		string(domain.CreateUser): ub.handleCreate,
		// string(domain.DeleteUser): ub.handleDelete,
	}
	return ub
}
