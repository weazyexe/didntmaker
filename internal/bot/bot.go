package bot

import (
	"log/slog"
	"time"

	"weazyexe.dev/didntmaker/internal/config"
	"weazyexe.dev/didntmaker/internal/handlers"
	"weazyexe.dev/didntmaker/internal/i18n"
	"weazyexe.dev/didntmaker/internal/repository"
	"weazyexe.dev/didntmaker/internal/service"

	"golang.org/x/time/rate"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	bot     *tele.Bot
	discord service.DiscordService
}

func New(cfg *config.Config, userRepo repository.UserRepository, postingRepo repository.PostingRepository, discordBindingRepo repository.DiscordBindingRepository) (*Bot, error) {
	return NewWithLang(cfg, userRepo, postingRepo, discordBindingRepo, i18n.Default())
}

func NewWithLang(cfg *config.Config, userRepo repository.UserRepository, postingRepo repository.PostingRepository, discordBindingRepo repository.DiscordBindingRepository, lang i18n.Lang) (*Bot, error) {
	slog.Info("creating bot instance", "lang", lang)

	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	messages := i18n.Get(lang)
	b, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("failed to create telebot", "error", err)
		return nil, err
	}

	b.Use(middleware.Recover())
	b.Use(newRateLimiter(rate.Limit(cfg.RateLimitPerSec), cfg.RateLimitBurst, messages).Middleware())

	userService := service.NewUserService(userRepo, postingRepo, cfg.DailyLimit)
	balanceService := service.NewBalanceService(userRepo, postingRepo, cfg.DailyLimit, cfg.SuperAdmin)
	betService := service.NewBetService(postingRepo, cfg.DailyLimit)

	h := handlers.New(b, userService, balanceService, betService, discordBindingRepo, messages)
	h.Register()

	result := &Bot{bot: b}

	// Initialize Discord service if token is provided
	if cfg.DiscordToken != "" && discordBindingRepo != nil {
		discordService, err := service.NewDiscordService(cfg.DiscordToken, discordBindingRepo, b, messages)
		if err != nil {
			slog.Error("failed to create discord service", "error", err)
		} else {
			result.discord = discordService
			slog.Info("discord service created")
		}
	}

	slog.Info("bot created", "username", b.Me.Username)
	return result, nil
}

func (b *Bot) Start() {
	if b.discord != nil {
		if err := b.discord.Start(); err != nil {
			slog.Error("failed to start discord service", "error", err)
		}
	}

	slog.Info("bot polling started")
	b.bot.Start()
}

func (b *Bot) Stop() {
	slog.Info("stopping bot")

	if b.discord != nil {
		b.discord.Stop()
	}

	b.bot.Stop()
}
