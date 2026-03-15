package event_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"poc-event-source/internal/application"
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/application/usecase/event"
	"poc-event-source/internal/domain"
)

// --- mocks ---

type mockEventRepo struct {
	createFn func(e *domain.EventSource) (*domain.EventSource, error)
}

func (m *mockEventRepo) CreateEvent(e *domain.EventSource) (*domain.EventSource, error) {
	return m.createFn(e)
}

type mockBroker struct {
	publishFn func(ctx context.Context, topic string, data []byte) error
}

func (m *mockBroker) Publish(ctx context.Context, topic string, data []byte) error {
	return m.publishFn(ctx, topic, data)
}
func (m *mockBroker) Subscribe(_ context.Context, _ string, _ func(context.Context, *application.Message)) (application.Subscription, error) {
	return nil, nil
}
func (m *mockBroker) QueueSubscribe(_ context.Context, _, _ string, _ func(context.Context, *application.Message)) (application.Subscription, error) {
	return nil, nil
}
func (m *mockBroker) Close() error { return nil }

// --- tests ---

func TestMainHandlerUseCase_Handler(t *testing.T) {
	ctx := context.Background()
	validEvent := dto.EventReqDTO{
		Type:    string(domain.CreateUser),
		Payload: datatypes.JSON(`{"username":"alice","password":"secret"}`),
	}

	tests := []struct {
		name      string
		repoErr   error
		brokerErr error
		wantErr   bool
	}{
		{
			name:    "saves event and publishes to broker successfully",
			wantErr: false,
		},
		{
			name:    "repo error does not publish to broker",
			repoErr: errors.New("db error"),
			wantErr: true,
		},
		{
			name:      "broker error returns error",
			brokerErr: errors.New("nats error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publishCalled := false

			repo := &mockEventRepo{
				createFn: func(e *domain.EventSource) (*domain.EventSource, error) {
					return e, tt.repoErr
				},
			}
			broker := &mockBroker{
				publishFn: func(_ context.Context, topic string, data []byte) error {
					publishCalled = true
					assert.Equal(t, "user", topic)
					assert.NotEmpty(t, data)
					return tt.brokerErr
				},
			}

			handler := event.NewMainHandler(repo, broker)
			err := handler.Handler(ctx, "user", validEvent)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, publishCalled, "broker.Publish must be called on happy path")
			}

			if tt.repoErr != nil {
				assert.False(t, publishCalled, "broker.Publish must not be called when repo fails")
			}
		})
	}
}
