package messaging_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"poc-event-source/internal/api/messaging"
	"poc-event-source/internal/application"
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/domain"
	"poc-event-source/internal/infrastructure/model"
)

// --- mocks ---

type mockSubscription struct{}

func (m *mockSubscription) Unsubscribe() error { return nil }

// syncBroker invokes the handler synchronously when Subscribe is called, simplifying tests
type syncBroker struct {
	msg *application.Message
}

func (b *syncBroker) Subscribe(ctx context.Context, _ string, handler func(context.Context, *application.Message)) (application.Subscription, error) {
	handler(ctx, b.msg)
	return &mockSubscription{}, nil
}
func (b *syncBroker) Publish(_ context.Context, _ string, _ []byte) error { return nil }
func (b *syncBroker) QueueSubscribe(_ context.Context, _, _ string, _ func(context.Context, *application.Message)) (application.Subscription, error) {
	return nil, nil
}
func (b *syncBroker) Close() error { return nil }

type mockUserRepo struct {
	createFn func(u *model.User) (*model.User, error)
}

func (m *mockUserRepo) CreateUser(u *model.User) (*model.User, error) {
	return m.createFn(u)
}

// --- helpers ---

func buildMessage(t *testing.T, eventType string, payload interface{}) (*application.Message, *bool) {
	t.Helper()
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)
	envelope, err := json.Marshal(dto.EventMessage{
		Type:    eventType,
		Payload: payloadBytes,
	})
	require.NoError(t, err)
	acked := false
	msg := &application.Message{
		Topic: "user",
		Data:  envelope,
		Ack:   func() error { acked = true; return nil },
	}
	return msg, &acked
}

// --- tests ---

func TestUserBroker_Subscribe_handleCreate(t *testing.T) {
	tests := []struct {
		name       string
		eventType  string
		payload    dto.CreateUserReqDTO
		wantCreate bool
		wantAck    bool
	}{
		{
			name:       "creates user successfully",
			eventType:  string(domain.CreateUser),
			payload:    dto.CreateUserReqDTO{Username: "alice", Password: "$2a$10$hashedpwd"},
			wantCreate: true,
			wantAck:    true,
		},
		{
			name:      "unknown event type — ack only, no user created",
			eventType: "UNKNOWN_EVENT",
			payload:   dto.CreateUserReqDTO{Username: "alice", Password: "$2a$10$hashedpwd"},
			wantAck:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created := false
			msg, acked := buildMessage(t, tt.eventType, tt.payload)

			broker := &syncBroker{msg: msg}

			repo := &mockUserRepo{
				createFn: func(u *model.User) (*model.User, error) {
					created = true
					if tt.wantCreate {
						assert.Equal(t, tt.payload.Username, u.Username)
						assert.Equal(t, tt.payload.Password, u.Password)
					}
					return u, nil
				},
			}

			userBroker := messaging.NewUserBroker(broker, repo)
			err := userBroker.Subscribe()

			assert.NoError(t, err)
			assert.Equal(t, tt.wantAck, *acked)
			assert.Equal(t, tt.wantCreate, created)
		})
	}
}

func TestUserBroker_Subscribe_handleCreate_dbError_noAck(t *testing.T) {
	msg, acked := buildMessage(t, string(domain.CreateUser), dto.CreateUserReqDTO{Username: "alice", Password: "$2a$10$hashedpwd"})

	broker := &syncBroker{msg: msg}
	repo := &mockUserRepo{
		createFn: func(_ *model.User) (*model.User, error) {
			return nil, errors.New("db error")
		},
	}

	userBroker := messaging.NewUserBroker(broker, repo)
	err := userBroker.Subscribe()

	assert.NoError(t, err)
	assert.False(t, *acked, "msg.Ack() must NOT be called when CreateUser fails")
}

func TestUserBroker_Subscribe_invalidJSON(t *testing.T) {
	acked := false
	msg := &application.Message{
		Topic: "user",
		Data:  []byte("not-json"),
		Ack:   func() error { acked = true; return nil },
	}

	broker := &syncBroker{msg: msg}
	repo := &mockUserRepo{
		createFn: func(_ *model.User) (*model.User, error) {
			t.Fatal("CreateUser must not be called with invalid JSON")
			return nil, nil
		},
	}

	userBroker := messaging.NewUserBroker(broker, repo)
	err := userBroker.Subscribe()

	assert.NoError(t, err)
	assert.True(t, acked, "msg.Ack() must be called even with invalid JSON")
}
