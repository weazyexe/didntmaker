package domain

import "time"

type User struct {
	ID         int64
	TelegramID int64
	ChatID     int64
	Username   string
	FirstName  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
