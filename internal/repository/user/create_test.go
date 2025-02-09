package user_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"poc-event-source/internal/infrastructure/model"
	"poc-event-source/internal/repository/testutils"
	"poc-event-source/internal/repository/user"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("Create user correctly", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		db, err := testutils.NewTestDatabase(ctx, &model.User{})
		assert.NoErrorf(t, err, "Error creating test database %v", err)

		tx := db.GormDB.Begin()
		defer tx.Rollback()

		repo := user.NewUserRepository(tx)

		userToBeCreated := &model.User{
			Password: "password",
			Username: "username",
		}

		res, err := repo.CreateUser(userToBeCreated)
		assert.NoError(t, err)

		var userFound model.User
		tx.First(&userFound, "username = ?", userToBeCreated.Username)

		assert.NotEmpty(t, res)
		assert.Equal(t, userFound.Username, res.Username)
		assert.Equal(t, userFound.Password, res.Password)
		assert.Equal(t, userFound.ID, res.ID)
	})
}
