package messaging

import (
	"poc-event-source/internal/api"
	"poc-event-source/internal/application"
)

type UserBroker struct {
	broker application.Broker
}

func NewUserBroker(broker application.Broker) api.Subscriber {
	return &UserBroker{broker: broker}
}
