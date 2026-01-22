package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotToken   string `env:"BOT_TOKEN" env-required:"true"`
	SuperAdmin string `env:"SUPER_ADMIN" env-default:""`
	DailyLimit int64  `env:"DAILY_LIMIT" env-default:"1000"`
	DBPath     string `env:"DB_PATH" env-default:"didntmaker.db"`
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
