package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"poc-event-source/internal/application/dto"
	"poc-event-source/internal/application/usecase/user"
	"poc-event-source/internal/domain"
)

// --- mock ---

type mockEventHandler struct {
	receivedTopic string
	receivedEvent dto.EventReqDTO
	returnErr     error
}

func (m *mockEventHandler) Handler(_ context.Context, topic string, event dto.EventReqDTO) error {
	m.receivedTopic = topic
	m.receivedEvent = event
	return m.returnErr
}

// --- tests ---

func TestCreateUserUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		input      dto.CreateUserReqDTO
		handlerErr error
		wantErr    bool
		wantTopic  string
		wantType   string
	}{
		{
			name:      "delegates to eventHandler with correct topic and type",
			input:     dto.CreateUserReqDTO{Username: "alice", Password: "secret"},
			wantTopic: "user",
			wantType:  string(domain.CreateUser),
			wantErr:   false,
		},
		{
			name:       "propagates eventHandler error",
			input:      dto.CreateUserReqDTO{Username: "alice", Password: "secret"},
			handlerErr: errors.New("handler error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockEventHandler{returnErr: tt.handlerErr}
			uc := user.NewCreateUserUseCase(mock)

			err := uc.Execute(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTopic, mock.receivedTopic)
			assert.Equal(t, tt.wantType, mock.receivedEvent.Type)
			assert.Contains(t, string(mock.receivedEvent.Payload), tt.input.Username)
		})
	}
}
