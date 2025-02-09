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
			name:     "senha válida deve gerar hash sem erros",
			password: "SenhaF0rte!2023",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "Erro inesperado")
				assert.NotEmpty(t, hash, "Hash gerado não pode ser vazio")

				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("SenhaF0rte!2023"))
				assert.NoError(t, err, "Hash não corresponde à senha")
			},
		},
		{
			name:     "hashs diferentes para mesma senha",
			password: "outraSenha123@",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "Erro inesperado")

				hash2, err := defaultCostUtils.HashPassword("outraSenha123@")
				assert.NoError(t, err, "Erro inesperado ao gerar segundo hash")
				assert.NotEqual(t, hash, hash2, "Hashs idênticos para a mesma senha - salting não está funcionando")
			},
		},
		{
			name:     "senha vazia deve ser tratada",
			password: "",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "Erro inesperado para senha vazia")
				assert.GreaterOrEqual(t, utf8.RuneCountInString(hash), 60, "Hash inválido para senha vazia")
			},
		},
		{
			name:     "senha muito longa",
			password: string(make([]byte, 100)),
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
			},
		},
		{
			name:     "caracteres especiais e unicode",
			password: "P@$$w0rd_áéíóú",
			setup:    func() application.PasswordUtil { return defaultCostUtils },
			assertFunc: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err, "Erro inesperado")
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("P@$$w0rd_áéíóú"))
				assert.NoError(t, err, "Hash não corresponde para caso de caracteres especiais")
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

	t.Run("teste de estresse com múltiplas requisições", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			_, err := defaultCostUtils.HashPassword("senhaTeste")
			assert.NoError(t, err, "Falha na iteração %d", i)
		}
	})
}
