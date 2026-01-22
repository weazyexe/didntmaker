package handlers

import tele "gopkg.in/telebot.v3"

// Start handles /start command
func (h *Handlers) Start(c tele.Context) error {
	logCommand(c, "/start")
	return c.Send(h.msg.Start)
}
