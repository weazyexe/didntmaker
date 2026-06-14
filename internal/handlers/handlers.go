package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"regexp"
	"time"

	"weazyexe.dev/didntmaker/internal/i18n"
	"weazyexe.dev/didntmaker/internal/repository"
	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

var deltaRe = regexp.MustCompile(`^([+-]?\d+)$`)

type Handlers struct {
	bot                *tele.Bot
	userService        service.UserService
	balanceService     service.BalanceService
	betService         service.BetService
	discordBindingRepo repository.DiscordBindingRepository
	msg                *i18n.Messages
}

func New(
	bot *tele.Bot,
	userService service.UserService,
	balanceService service.BalanceService,
	betService service.BetService,
	discordBindingRepository repository.DiscordBindingRepository,
	msg *i18n.Messages,
) *Handlers {
	return &Handlers{
		bot:                bot,
		userService:        userService,
		balanceService:     balanceService,
		betService:         betService,
		discordBindingRepo: discordBindingRepository,
		msg:                msg,
	}
}

func (h *Handlers) Register() {
	h.bot.Handle("/start", h.Start)
	h.bot.Handle("/help", h.Help)
	h.bot.Handle("/me", h.Me)
	h.bot.Handle("/stats", h.Stats)
	h.bot.Handle("/balances", h.Balances)
	h.bot.Handle("/bet", h.Bet)
	h.bot.Handle("/bet_stats", h.BetStats)
	h.bot.Handle("/add", h.Add)
	if h.discordBindingRepo != nil {
		h.bot.Handle("/discord_bind", h.DiscordBind)
		h.bot.Handle("/discord_unbind", h.DiscordUnbind)
	}

	h.bot.Handle(tele.OnText, h.Reply)

	slog.Info("handlers registered")
}

func logCommand(c tele.Context, command string) func() {
	reqID := newRequestID()
	start := time.Now()
	slog.Info("command received",
		"req_id", reqID,
		"command", command,
		"user_id", c.Sender().ID,
		"username", c.Sender().Username,
		"chat_id", c.Chat().ID,
		"payload", c.Message().Payload,
	)
	return func() {
		slog.Info("command handled", "req_id", reqID, "command", command, "duration", time.Since(start))
	}
}

func newRequestID() string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
