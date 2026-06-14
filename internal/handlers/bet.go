package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"weazyexe.dev/didntmaker/internal/domain"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Bet(c tele.Context) error {
	defer logCommand(c, "/bet")()
	ctx := context.Background()
	chatID := c.Chat().ID
	telegramID := c.Sender().ID

	_, err := h.userService.GetOrCreate(ctx, chatID, telegramID, c.Sender().Username, c.Sender().FirstName)
	if err != nil {
		return c.Send(h.msg.BetNotRegistered)
	}

	if err := h.betService.CanBet(ctx, chatID, telegramID); err != nil {
		if errors.Is(err, domain.ErrBetAlreadyUsed) {
			return c.Send(h.msg.BetAlreadyUsed)
		}
		return c.Send(h.msg.BetError)
	}

	msg, err := tele.Cube.Send(h.bot, c.Chat(), nil)
	if err != nil {
		return c.Send(h.msg.BetDiceError)
	}

	diceValue := msg.Dice.Value

	result, err := h.betService.ApplyResult(ctx, chatID, telegramID, diceValue)
	if err != nil {
		return c.Send(h.msg.BetResultError)
	}

	time.Sleep(4 * time.Second)

	if result.Won {
		return c.Send(fmt.Sprintf(h.msg.BetWin, result.DailyLimit))
	}
	return c.Send(fmt.Sprintf(h.msg.BetLose, result.DailyLimit))
}
