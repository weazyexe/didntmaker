package handlers

import (
	"context"
	"fmt"
	"strings"

	"weazyexe.dev/didntmaker/internal/domain"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Stats(c tele.Context) error {
	defer logCommand(c, "/stats")()
	users, err := h.userService.GetLeaderboard(context.Background(), c.Chat().ID)
	return h.sendLeaderboard(c, h.msg.StatsHeader, users, err)
}

func (h *Handlers) StatsWeek(c tele.Context) error {
	defer logCommand(c, "/stats_week")()
	users, err := h.userService.GetWeeklyLeaderboard(context.Background(), c.Chat().ID)
	return h.sendLeaderboard(c, h.msg.StatsWeekHeader, users, err)
}

func (h *Handlers) StatsMonth(c tele.Context) error {
	defer logCommand(c, "/stats_month")()
	users, err := h.userService.GetMonthlyLeaderboard(context.Background(), c.Chat().ID)
	return h.sendLeaderboard(c, h.msg.StatsMonthHeader, users, err)
}

func (h *Handlers) sendLeaderboard(c tele.Context, header string, users []domain.LeaderboardEntry, err error) error {
	if err != nil {
		return c.Send(h.msg.StatsError)
	}

	if len(users) == 0 {
		return c.Send(h.msg.StatsEmpty)
	}

	var sb strings.Builder
	sb.WriteString(header)

	for i, user := range users {
		var prefix string
		if i < len(h.msg.StatsMedals) {
			prefix = h.msg.StatsMedals[i]
		} else {
			prefix = fmt.Sprintf("%d.", i+1)
		}

		displayName := user.FirstName
		if user.Username != "" {
			displayName = "@" + user.Username
		}

		sb.WriteString(fmt.Sprintf(h.msg.StatsEntryFmt, prefix, displayName, user.Score))
	}

	return c.Send(sb.String())
}
