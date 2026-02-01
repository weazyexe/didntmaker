package handlers

import (
	"fmt"
	"regexp"

	tele "gopkg.in/telebot.v3"
	"weazyexe.dev/didntmaker/internal/models"
)

var guildIDRe = regexp.MustCompile(`^\d+$`)

func (h *Handlers) DiscordBind(c tele.Context) error {
	logCommand(c, "/discord_bind")

	guildID := c.Message().Payload
	if guildID == "" {
		return c.Send(h.msg.DiscordBindUsage)
	}

	if !guildIDRe.MatchString(guildID) {
		return c.Send(h.msg.DiscordBindInvalidID)
	}

	chatID := c.Chat().ID

	exists, err := h.discordBindingRepo.Exists(chatID, guildID)
	if err != nil {
		return c.Send(h.msg.DiscordBindError)
	}
	if exists {
		return c.Send(h.msg.DiscordBindAlreadyBound)
	}

	binding := &models.DiscordBinding{
		ChatID:  chatID,
		GuildID: guildID,
	}
	if err := h.discordBindingRepo.Create(binding); err != nil {
		return c.Send(h.msg.DiscordBindError)
	}

	return c.Send(fmt.Sprintf(h.msg.DiscordBindSuccess, guildID))
}

func (h *Handlers) DiscordUnbind(c tele.Context) error {
	logCommand(c, "/discord_unbind")

	guildID := c.Message().Payload
	if guildID == "" {
		return c.Send(h.msg.DiscordUnbindUsage)
	}

	chatID := c.Chat().ID

	if err := h.discordBindingRepo.Delete(chatID, guildID); err != nil {
		return c.Send(h.msg.DiscordUnbindError)
	}

	return c.Send(fmt.Sprintf(h.msg.DiscordUnbindSuccess, guildID))
}
