package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TelegramID   int64  `gorm:"index:idx_telegram_chat,unique"`
	ChatID       int64  `gorm:"index:idx_telegram_chat,unique"`
	Username     string
	FirstName    string
	Balance      int64     `gorm:"default:0"`
	DailyGiven   int64     `gorm:"default:0"`
	DailyResetAt time.Time `gorm:"autoCreateTime"`
	LastBetAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
