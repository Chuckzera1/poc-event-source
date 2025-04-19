package application

import (
	"context"
)

type Message struct {
	Topic string
	Data  []byte

	Ack func() error
}

type Subscription interface {
	Unsubscribe() error
}

type Broker interface {
	Publish(ctx context.Context, topic string, data []byte) error

	Subscribe(ctx context.Context, topic string, handler func(context.Context, *Message)) (Subscription, error)

	QueueSubscribe(ctx context.Context, topic, queue string, handler func(context.Context, *Message)) (Subscription, error)

	Close() error
}
