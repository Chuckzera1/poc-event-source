package utils

import "poc-event-source/internal/application"

type passwordUtils struct {
	cost int
}

func NewPasswordUtils(cost int) application.PasswordUtil {
	return &passwordUtils{
		cost: cost,
	}
}
