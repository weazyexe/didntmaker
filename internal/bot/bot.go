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
	bot *tele.Bot
}

func New(cfg *config.Config, userRepo *repository.UserRepository) (*Bot, error) {
	return NewWithLang(cfg, userRepo, i18n.Default())
}

func NewWithLang(cfg *config.Config, userRepo *repository.UserRepository, lang i18n.Lang) (*Bot, error) {
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
	userSvc := service.NewUserService(userRepo)
	balanceSvc := service.NewBalanceService(userRepo, cfg.SuperAdmin)
	betSvc := service.NewBetService(userRepo)

	// Create and register handlers
	h := handlers.New(b, userSvc, balanceSvc, betSvc, i18n.Get(lang))
	h.Register()

	slog.Info("bot created", "username", b.Me.Username)
	return &Bot{bot: b}, nil
}

func (b *Bot) Start() {
	slog.Info("bot polling started")
	b.bot.Start()
}
