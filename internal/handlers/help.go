package handlers

import tele "gopkg.in/telebot.v3"

// Help handles /help command
func (h *Handlers) Help(c tele.Context) error {
	logCommand(c, "/help")
	return c.Send(h.msg.Help)
}
