package handlers

import (
	"context"
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
	h.bot.Use(h.ensureRegistered)

	h.bot.Handle("/info", h.Info)
	h.bot.Handle("/stats", h.Stats)
	h.bot.Handle("/stats_week", h.StatsWeek)
	h.bot.Handle("/stats_month", h.StatsMonth)
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

// ensureRegistered registers (or refreshes) the sender on every interaction, so
// anyone who writes in the chat is on the books — no explicit /register needed.
func (h *Handlers) ensureRegistered(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if s := c.Sender(); s != nil && !s.IsBot && c.Chat() != nil {
			if _, err := h.userService.GetOrCreate(context.Background(), c.Chat().ID, s.ID, s.Username, s.FirstName); err != nil {
				slog.Warn("auto-register failed", "user_id", s.ID, "chat_id", c.Chat().ID, "error", err)
			}
		}
		return next(c)
	}
}

// displayName prefers @username, falling back to first name.
func displayName(username, firstName string) string {
	if username != "" {
		return "@" + username
	}
	return firstName
}
