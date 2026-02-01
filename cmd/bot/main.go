package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"weazyexe.dev/didntmaker/internal/bot"
	"weazyexe.dev/didntmaker/internal/config"
	"weazyexe.dev/didntmaker/internal/database"
	"weazyexe.dev/didntmaker/internal/repository"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting application")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded")

	db, err := database.Init(cfg.DBPath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	slog.Info("database initialized")

	userRepo := repository.NewUserRepository(db, cfg.DailyLimit)
	txRepo := repository.NewTransactionRepository(db)

	var discordBindingRepo repository.DiscordBindingRepository
	if cfg.DiscordToken != "" {
		discordBindingRepo = repository.NewDiscordBindingRepository(db)
		slog.Info("discord integration enabled")
	}

	b, err := bot.New(cfg, userRepo, txRepo, discordBindingRepo)
	if err != nil {
		slog.Error("failed to create bot", "error", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("bot starting")
		b.Start()
	}()

	<-sig
	slog.Info("shutdown signal received")

	b.Stop()

	if err := database.Close(db); err != nil {
		slog.Error("failed to close database", "error", err)
	}

	slog.Info("application stopped")
}
