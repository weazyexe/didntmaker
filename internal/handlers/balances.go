package handlers

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) Balances(c tele.Context) error {
	logCommand(c, "/balances")

	balances, err := h.balanceService.GetDailyBalances(c.Chat().ID)
	if err != nil {
		return c.Send(h.msg.BalancesError)
	}

	if len(balances) == 0 {
		return c.Send(h.msg.BalancesEmpty)
	}

	var sb strings.Builder
	sb.WriteString(h.msg.BalancesHeader)

	for _, b := range balances {
		displayName := b.User.FirstName
		if b.User.Username != "" {
			displayName = "@" + b.User.Username
		}

		status := ""
		if b.Remaining == 0 {
			status = h.msg.BalancesEmpty_
		} else if b.Remaining == b.DailyLimit {
			status = h.msg.BalancesFull
		}

		sb.WriteString(fmt.Sprintf(h.msg.BalancesEntry, displayName, b.Remaining, b.DailyLimit, status))
	}

	return c.Send(sb.String())
}
