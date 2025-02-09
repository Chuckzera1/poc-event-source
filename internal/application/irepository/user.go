package irepository

import "poc-event-source/internal/infrastructure/model"

type CreateUserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
}

type UserRepository interface {
	CreateUserRepository
}
