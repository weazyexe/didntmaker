package handlers

import (
	"log/slog"
	"regexp"

	"weazyexe.dev/didntmaker/internal/i18n"
	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

var deltaRe = regexp.MustCompile(`^([+-]?\d+)$`)

type Handlers struct {
	bot            *tele.Bot
	userService    service.UserService
	balanceService service.BalanceService
	betService     service.BetService
	msg            *i18n.Messages
}

func New(
	bot *tele.Bot,
	userSvc service.UserService,
	balanceSvc service.BalanceService,
	betSvc service.BetService,
	msg *i18n.Messages,
) *Handlers {
	return &Handlers{
		bot:            bot,
		userService:    userSvc,
		balanceService: balanceSvc,
		betService:     betSvc,
		msg:            msg,
	}
}

func (h *Handlers) Register() {
	h.bot.Handle("/start", h.Start)
	h.bot.Handle("/help", h.Help)
	h.bot.Handle("/me", h.Me)
	h.bot.Handle("/stats", h.Stats)
	h.bot.Handle("/balances", h.Balances)
	h.bot.Handle("/bet", h.Bet)
	h.bot.Handle("/add", h.Add)
	h.bot.Handle(tele.OnText, h.Reply)

	slog.Info("handlers registered")
}

func logCommand(c tele.Context, command string) {
	slog.Info("command received",
		"command", command,
		"user_id", c.Sender().ID,
		"username", c.Sender().Username,
		"chat_id", c.Chat().ID,
		"payload", c.Message().Payload,
	)
}
