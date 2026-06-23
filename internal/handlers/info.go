package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"weazyexe.dev/didntmaker/internal/domain"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Info(c tele.Context) error {
	defer logCommand(c, "/info")()

	// no reply -> own stats; reply -> the replied-to user's stats
	target := c.Sender()
	if reply := c.Message().ReplyTo; reply != nil && reply.Sender != nil {
		if reply.Sender.ID == h.bot.Me.ID {
			return c.Send(h.msg.InfoOnBot)
		}
		target = reply.Sender
	}
	return h.sendStats(c, target)
}

func (h *Handlers) sendStats(c tele.Context, target *tele.User) error {
	stats, err := h.userService.GetStats(
		context.Background(),
		c.Chat().ID,
		target.ID,
		target.Username,
		target.FirstName,
	)
	if err != nil {
		return c.Send(h.msg.MeError)
	}

	var b strings.Builder

	fmt.Fprintf(&b, h.msg.MeHeader, displayName(stats.User.Username, stats.User.FirstName))
	fmt.Fprintf(&b, h.msg.MeScore, formatNum(stats.Score))
	fmt.Fprintf(&b, h.msg.MeWeek, signedTrend(stats.WeekDelta))
	fmt.Fprintf(&b, h.msg.MeMonth, signedTrend(stats.MonthDelta))

	if stats.WorstDayMinus > 0 || stats.BestDayPlus > 0 {
		b.WriteString("\n")
	}
	if stats.WorstDayMinus > 0 {
		fmt.Fprintf(&b, h.msg.MeWorstDay, formatNum(stats.WorstDayMinus))
	}
	if stats.BestDayPlus > 0 {
		fmt.Fprintf(&b, h.msg.MeBestDay, formatNum(stats.BestDayPlus))
	}

	fmt.Fprintf(&b, h.msg.MeLimit, formatNum(stats.DailyRemaining), formatNum(stats.DailyLimit))
	fmt.Fprintf(&b, h.msg.MeBets, stats.Won, stats.Lost)
	if stats.BetAvailable {
		b.WriteString(h.msg.MeBetReady)
	} else {
		b.WriteString(h.msg.MeBetUsed)
	}

	if stats.Fan != nil || stats.Hater != nil || stats.Favorite != nil || stats.Victim != nil {
		b.WriteString("\n")
		writeCounterparty(&b, h.msg.MeFan, stats.Fan)
		writeCounterparty(&b, h.msg.MeHater, stats.Hater)
		writeCounterparty(&b, h.msg.MeFavorite, stats.Favorite)
		writeCounterparty(&b, h.msg.MeVictim, stats.Victim)
	}

	return c.Send(b.String())
}

func writeCounterparty(b *strings.Builder, format string, cp *domain.Counterparty) {
	if cp == nil {
		return
	}
	fmt.Fprintf(b, format, displayName(cp.Username, cp.FirstName), formatNum(cp.Amount))
}

// signedTrend renders a period delta as "+320 ▲" / "−180 ▼" / "0 ➖".
func signedTrend(delta int64) string {
	switch {
	case delta > 0:
		return "+" + formatNum(delta) + " ▲"
	case delta < 0:
		return "−" + formatNum(-delta) + " ▼"
	default:
		return "0 ➖"
	}
}

// formatNum groups thousands with a space: 1240 -> "1 240". Handles negatives.
func formatNum(n int64) string {
	s := strconv.FormatInt(n, 10)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}

	var out strings.Builder
	for i, d := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out.WriteRune(' ') // no-break space, keeps "1 240" from wrapping
		}
		out.WriteRune(d)
	}

	if neg {
		return "−" + out.String()
	}
	return out.String()
}
