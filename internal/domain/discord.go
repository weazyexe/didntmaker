package domain

import "time"

type DiscordBinding struct {
	ChatID    int64
	GuildID   string
	CreatedAt time.Time
}
