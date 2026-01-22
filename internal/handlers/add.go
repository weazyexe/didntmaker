package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Add(c tele.Context) error {
	logCommand(c, "/add")

	args := c.Message().Payload
	if args == "" {
		return c.Send(h.msg.AddUsage)
	}

	re := regexp.MustCompile(`@(\w+)\s+([+-]?\d+)`)
	matches := re.FindStringSubmatch(args)

	if len(matches) != 3 {
		return c.Send(h.msg.AddFormatError)
	}

	username := matches[1]
	delta, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return c.Send(h.msg.AddNumberError)
	}

	result, err := h.balanceService.AdjustDailyLimit(
		c.Chat().ID,
		c.Sender().Username,
		username,
		delta,
	)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			return nil // Silently ignore
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Send(fmt.Sprintf(h.msg.AddNotFound, username))
		}
		return c.Send(h.msg.AddError)
	}

	sign := ""
	if delta > 0 {
		sign = "+"
	}

	return c.Send(fmt.Sprintf(h.msg.AddSuccess, username, result.OldRemaining, result.NewRemaining, sign, delta))
}
