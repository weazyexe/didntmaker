package models

import (
	"time"
)

type DiscordBinding struct {
	ID        uint   `gorm:"primaryKey"`
	ChatID    int64  `gorm:"index:idx_discord_binding,unique"`
	GuildID   string `gorm:"index:idx_discord_binding,unique;index:idx_guild"`
	CreatedAt time.Time
}
