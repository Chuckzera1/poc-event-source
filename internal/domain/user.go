package domain

import "time"

type User struct {
	ID uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Username string
	Password string
}
