package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Reply(c tele.Context) error {
	reply := c.Message().ReplyTo
	if reply == nil {
		return nil
	}

	text := strings.TrimSpace(c.Message().Text)
	match := deltaRe.FindStringSubmatch(text)
	if match == nil {
		return nil
	}

	delta, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil || delta == 0 {
		return nil
	}

	logCommand(c, "reply")

	target := reply.Sender
	if target == nil {
		return c.Send(h.msg.ReplyUnknownTarget)
	}

	chatID := c.Chat().ID
	senderID := c.Sender().ID

	// Reply to bot -> transfer to all
	if target.ID == h.bot.Me.ID {
		return h.replyToAll(c, chatID, senderID, delta)
	}

	result, err := h.balanceService.Transfer(chatID, senderID, target.ID, target.Username, delta)
	if err != nil {
		return h.handleTransferError(c, err)
	}

	displayName := result.Target.FirstName
	if result.Target.Username != "" {
		displayName = "@" + result.Target.Username
	}

	if delta < 0 {
		return c.Send(fmt.Sprintf(h.msg.ReplySuccessNegative, displayName, delta))
	}
	return c.Send(fmt.Sprintf(h.msg.ReplySuccessPositive, displayName, delta))
}

func (h *Handlers) replyToAll(c tele.Context, chatID, senderID, delta int64) error {
	result, err := h.balanceService.TransferToAll(chatID, senderID, delta)
	if err != nil {
		if errors.Is(err, service.ErrNoUsersInChat) {
			return c.Send(h.msg.ReplyAllNoUsers)
		}
		if errors.Is(err, service.ErrInsufficientLimit) {
			return c.Send(fmt.Sprintf(h.msg.ReplyAllNotEnough, result.TotalCost, 0))
		}
		return c.Send(h.msg.ReplyAllError)
	}

	if delta < 0 {
		return c.Send(fmt.Sprintf(h.msg.ReplyAllSuccessNeg, delta))
	}
	return c.Send(fmt.Sprintf(h.msg.ReplyAllSuccessPos, delta))
}

func (h *Handlers) handleTransferError(c tele.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrSelfTransfer):
		return c.Send(h.msg.ReplySelfError)
	case errors.Is(err, service.ErrTransactionLimit):
		return c.Send(h.msg.ReplyLimitExceeded)
	case errors.Is(err, service.ErrUserNotFound):
		return c.Send(h.msg.ReplyTargetNotFound)
	case errors.Is(err, service.ErrInsufficientLimit):
		return c.Send(fmt.Sprintf(h.msg.ReplyNotEnough, 0, h.balanceService.DailyLimit()))
	default:
		return c.Send(h.msg.ReplyError)
	}
}
