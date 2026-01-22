package handlers

import (
	"errors"
	"fmt"
	"time"

	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Bet(c tele.Context) error {
	logCommand(c, "/bet")
	chatID := c.Chat().ID
	telegramID := c.Sender().ID

	_, err := h.userService.GetOrCreate(chatID, telegramID, c.Sender().Username, c.Sender().FirstName)
	if err != nil {
		return c.Send(h.msg.BetNotRegistered)
	}

	if err := h.betService.CanBet(chatID, telegramID); err != nil {
		if errors.Is(err, service.ErrBetAlreadyUsed) {
			return c.Send(h.msg.BetAlreadyUsed)
		}
		return c.Send(h.msg.BetError)
	}

	msg, err := tele.Cube.Send(h.bot, c.Chat(), nil)
	if err != nil {
		return c.Send(h.msg.BetDiceError)
	}

	diceValue := msg.Dice.Value

	result, err := h.betService.ApplyResult(chatID, telegramID, diceValue)
	if err != nil {
		return c.Send(h.msg.BetResultError)
	}

	time.Sleep(4 * time.Second)

	if result.Won {
		return c.Send(fmt.Sprintf(h.msg.BetWin, result.DiceValue, result.DailyLimit))
	}
	return c.Send(fmt.Sprintf(h.msg.BetLose, result.DiceValue, result.DailyLimit))
}
