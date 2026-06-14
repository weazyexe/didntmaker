package handlers

import (
	"context"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) BetStats(c tele.Context) error {
	defer logCommand(c, "/bet_stats")()

	stats, err := h.betService.ChatBetStats(context.Background(), c.Chat().ID)
	if err != nil {
		return c.Send(h.msg.BetStatsError)
	}

	if len(stats) == 0 {
		return c.Send(h.msg.BetStatsEmpty)
	}

	var sb strings.Builder
	sb.WriteString(h.msg.BetStatsHeader)

	for _, s := range stats {
		displayName := s.FirstName
		if s.Username != "" {
			displayName = "@" + s.Username
		}
		sb.WriteString(fmt.Sprintf(h.msg.BetStatsEntry, displayName, s.Won, s.Lost))
	}

	return c.Send(sb.String())
}
