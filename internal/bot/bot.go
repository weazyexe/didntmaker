package bot

import (
	"log/slog"
	"time"

	"weazyexe.dev/didntmaker/internal/config"
	"weazyexe.dev/didntmaker/internal/handlers"
	"weazyexe.dev/didntmaker/internal/i18n"
	"weazyexe.dev/didntmaker/internal/repository"
	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	bot     *tele.Bot
	discord service.DiscordService
}

func New(cfg *config.Config, userRepo repository.UserRepository, txRepo repository.TransactionRepository, discordBindingRepo repository.DiscordBindingRepository) (*Bot, error) {
	return NewWithLang(cfg, userRepo, txRepo, discordBindingRepo, i18n.Default())
}

func NewWithLang(cfg *config.Config, userRepo repository.UserRepository, txRepo repository.TransactionRepository, discordBindingRepo repository.DiscordBindingRepository, lang i18n.Lang) (*Bot, error) {
	slog.Info("creating bot instance", "lang", lang)

	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("failed to create telebot", "error", err)
		return nil, err
	}

	// Create services
	txSvc := service.NewTransactionService(txRepo)
	userSvc := service.NewUserService(userRepo)
	balanceSvc := service.NewBalanceService(userRepo, txSvc, cfg.SuperAdmin)
	betSvc := service.NewBetService(userRepo, txSvc)

	// Create and register handlers
	h := handlers.New(b, userSvc, balanceSvc, betSvc, txSvc, discordBindingRepo, i18n.Get(lang))
	h.Register()

	result := &Bot{bot: b}

	// Initialize Discord service if token is provided
	if cfg.DiscordToken != "" && discordBindingRepo != nil {
		discordSvc, err := service.NewDiscordService(cfg.DiscordToken, discordBindingRepo, b, i18n.Get(lang))
		if err != nil {
			slog.Error("failed to create discord service", "error", err)
		} else {
			result.discord = discordSvc
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
