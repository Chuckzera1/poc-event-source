package utils

import "golang.org/x/crypto/bcrypt"

func (p *passwordUtils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)

	return string(bytes), err
}

func (p *passwordUtils) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
