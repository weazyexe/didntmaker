package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotToken        string   `env:"BOT_TOKEN" env-required:"true"`
	SuperAdmin      []string `env:"SUPER_ADMIN" env-default:""`
	DailyLimit      int64    `env:"DAILY_LIMIT" env-default:"1000"`
	DBPath          string   `env:"DB_PATH" env-default:"didntmaker.db"`
	DiscordToken    string   `env:"DISCORD_TOKEN" env-default:""`
	RateLimitPerSec float64  `env:"RATE_LIMIT_PER_SEC" env-default:"1"`
	RateLimitBurst  int      `env:"RATE_LIMIT_BURST" env-default:"5"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		slog.Warn("no .env file found, reading from environment variables")
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
