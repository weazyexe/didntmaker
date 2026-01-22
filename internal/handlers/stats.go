package handlers

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Stats(c tele.Context) error {
	logCommand(c, "/stats")

	users, err := h.userService.GetLeaderboard(c.Chat().ID)
	if err != nil {
		return c.Send(h.msg.StatsError)
	}

	if len(users) == 0 {
		return c.Send(h.msg.StatsEmpty)
	}

	var sb strings.Builder
	sb.WriteString(h.msg.StatsHeader)

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

		sb.WriteString(fmt.Sprintf(h.msg.StatsEntryFmt, prefix, displayName, user.Balance))
	}

	return c.Send(sb.String())
}
