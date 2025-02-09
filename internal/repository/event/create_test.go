package event_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"poc-event-source/internal/infrastructure/model"
	"poc-event-source/internal/repository/event"
	"poc-event-source/internal/repository/testutils"
	"testing"
	"time"
)

func TestEventRepository_CreateEvent(t *testing.T) {
	ctx := context.Background()

	db, err := testutils.NewTestDatabase(ctx, &model.EventSource{})

	assert.NoError(t, err)

	now := time.Now().UTC()

	tests := []struct {
		name        string
		event       *model.EventSource
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully create event with all fields",
			event: &model.EventSource{
				AggregateID: uuid.NewString(),
				Type:        "UserCreated",
				Payload:     []byte(`{"name": "John"}`),
			},
			wantErr: false,
		},
		{
			name: "successfully create event with minimum required fields",
			event: &model.EventSource{
				AggregateID: uuid.NewString(),
				Type:        "UserUpdated",
				Payload:     []byte(`{}`),
			},
			wantErr: false,
		},
		{
			name: "fail when missing event type",
			event: &model.EventSource{
				AggregateID: uuid.NewString(),
				Payload:     []byte(`{"name": "John"}`),
			},
			wantErr:     true,
			errContains: "not-null constraint",
		},
		{
			name: "fail when missing payload",
			event: &model.EventSource{
				AggregateID: uuid.NewString(),
				Type:        "UserDeleted",
			},
			wantErr:     true,
			errContains: "not-null constraint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tx := db.GormDB.Begin()
			repo := event.NewEventRepository(tx)
			defer tx.Rollback()

			createdEvent, err2 := repo.CreateEvent(tt.event)

			if tt.wantErr {
				assert.Error(t, err2)
				if tt.errContains != "" {
					assert.Contains(t, err2.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err2)
			assert.NotEmpty(t, createdEvent.ID)
			assert.Equal(t, tt.event.AggregateID, createdEvent.AggregateID)
			assert.Equal(t, tt.event.Type, createdEvent.Type)
			assert.Equal(t, tt.event.Payload, createdEvent.Payload)
			assert.Equal(t, 1, createdEvent.Version, "Version should default to 1")
			assert.WithinDuration(t, now, createdEvent.CreatedAt, 5*time.Second, "CreatedAt should be set to current time")

			var dbEvent model.EventSource
			result := tx.First(&dbEvent, "id = ?", createdEvent.ID)
			assert.NoError(t, result.Error, "should find created event in database")
			assert.Equal(t, createdEvent.ID, dbEvent.ID)
			assert.Equal(t, createdEvent.AggregateID, dbEvent.AggregateID)
			assert.Equal(t, createdEvent.Type, dbEvent.Type)
			assert.Equal(t, createdEvent.Payload, dbEvent.Payload)
			assert.Equal(t, createdEvent.Version, dbEvent.Version)
			assert.WithinDuration(t, createdEvent.CreatedAt, dbEvent.CreatedAt, 5*time.Second)

			if tt.name == "fail with duplicate ID" {
				_, err := repo.CreateEvent(tt.event)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "duplicate key")
			}
		})
	}
}
