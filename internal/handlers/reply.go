package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"weazyexe.dev/didntmaker/internal/domain"

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

	defer logCommand(c, "reply")()

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

	result, err := h.balanceService.Transfer(context.Background(), chatID, senderID, target.ID, target.Username, delta)
	if err != nil {
		return h.handleTransferError(c, err, result)
	}

	name := displayName(result.Target.Username, result.Target.FirstName)

	if delta < 0 {
		return c.Send(fmt.Sprintf(h.msg.ReplySuccessNegative, name, delta))
	}
	return c.Send(fmt.Sprintf(h.msg.ReplySuccessPositive, name, delta))
}

func (h *Handlers) replyToAll(c tele.Context, chatID, senderID, delta int64) error {
	result, err := h.balanceService.TransferToAll(context.Background(), chatID, senderID, delta)
	if err != nil {
		if errors.Is(err, domain.ErrNoUsersInChat) {
			return c.Send(h.msg.ReplyAllNoUsers)
		}
		if errors.Is(err, domain.ErrInsufficientLimit) {
			return c.Send(fmt.Sprintf(h.msg.ReplyAllNotEnough, result.TotalCost, result.Remaining))
		}
		return c.Send(h.msg.ReplyAllError)
	}

	if delta < 0 {
		return c.Send(fmt.Sprintf(h.msg.ReplyAllSuccessNeg, delta))
	}
	return c.Send(fmt.Sprintf(h.msg.ReplyAllSuccessPos, delta))
}

func (h *Handlers) handleTransferError(c tele.Context, err error, result *domain.TransferResult) error {
	switch {
	case errors.Is(err, domain.ErrSelfTransfer):
		return c.Send(h.msg.ReplySelfError)
	case errors.Is(err, domain.ErrUserNotFound):
		return c.Send(h.msg.ReplyTargetNotFound)
	case errors.Is(err, domain.ErrInsufficientLimit):
		remaining := int64(0)
		if result != nil {
			remaining = result.Remaining
		}
		return c.Send(fmt.Sprintf(h.msg.ReplyNotEnough, remaining, h.balanceService.DailyLimit()))
	default:
		return c.Send(h.msg.ReplyError)
	}
}
