package user

import "poc-event-source/internal/infrastructure/model"

func (u userRepository) CreateUser(user *model.User) (*model.User, error) {
	err := u.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
