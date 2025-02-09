package application

type PasswordUtil interface {
	HashPassword(password string) (string, error)
}
