package user

import (
	"gorm.io/gorm"
	"poc-event-source/internal/application/irepository"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) irepository.UserRepository {
	return &userRepository{db: db}
}
