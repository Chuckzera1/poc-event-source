package utils_test

import (
	"poc-event-source/internal/application"
	"poc-event-source/internal/application/utils"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	defaultCostUtils := utils.NewPasswordUtils(bcrypt.DefaultCost)

	tests := []struct {
		name       string
		password   string
		setup      func() application.PasswordUtil
		assertFunc func(t *testing.T, hash string, err error)
	}{
		{
			name:     "valid password should generate hash without errors",
			password: "StrongP@ss!2023",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "unexpected error")
				assert.NotEmpty(t, hash, "generated hash must not be empty")

				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("StrongP@ss!2023"))
				assert.NoError(t, err, "hash does not match password")
			},
		},
		{
			name:     "different hashes for same password",
			password: "anotherPassword123@",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "unexpected error")

				hash2, err := defaultCostUtils.HashPassword("anotherPassword123@")
				assert.NoError(t, err, "unexpected error generating second hash")
				assert.NotEqual(t, hash, hash2, "identical hashes for same password — salting is not working")
			},
		},
		{
			name:     "empty password should be handled",
			password: "",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "unexpected error for empty password")
				assert.GreaterOrEqual(t, utf8.RuneCountInString(hash), 60, "invalid hash for empty password")
			},
		},
		{
			name:     "password too long",
			password: string(make([]byte, 100)),
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
			},
		},
		{
			name:     "special characters and unicode",
			password: "P@$$w0rd_áéíóú",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "unexpected error")
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("P@$$w0rd_áéíóú"))
				assert.NoError(t, err, "hash does not match for special characters case")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passwordUtils := tt.setup()
			hash, err := passwordUtils.HashPassword(tt.password)
			tt.assertFunc(t, hash, err)
		})
	}

	t.Run("stress test with multiple requests", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			_, err := defaultCostUtils.HashPassword("testPassword")
			assert.NoError(t, err, "failed on iteration %d", i)
		}
	})
}
