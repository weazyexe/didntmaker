package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Me(c tele.Context) error {
	logCommand(c, "/me")

	stats, err := h.userService.GetStats(
		c.Chat().ID,
		c.Sender().ID,
		c.Sender().Username,
		c.Sender().FirstName,
	)
	if err != nil {
		return c.Send(h.msg.MeError)
	}

	return c.Send(fmt.Sprintf(h.msg.MeStats, stats.User.Balance, stats.DailyRemaining, stats.DailyLimit))
}
