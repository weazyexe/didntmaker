package handlers

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Me(c tele.Context) error {
	defer logCommand(c, "/me")()

	stats, err := h.userService.GetStats(
		context.Background(),
		c.Chat().ID,
		c.Sender().ID,
		c.Sender().Username,
		c.Sender().FirstName,
	)
	if err != nil {
		return c.Send(h.msg.MeError)
	}

	msg := fmt.Sprintf(h.msg.MeStats, stats.Score, stats.DailyRemaining, stats.DailyLimit) +
		fmt.Sprintf(h.msg.MeBets, stats.Won, stats.Lost)
	return c.Send(msg)
}
