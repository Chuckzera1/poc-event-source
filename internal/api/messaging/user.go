package messaging

import (
	"context"
	"fmt"
	"poc-event-source/internal/application"
)

func (ub *UserBroker) Subscribe() error {
	ctx := context.Background()
	ub.broker.Subscribe(ctx, "user", func(ctx context.Context, msg *application.Message) {
		fmt.Println(msg)
	})
	return nil
}
