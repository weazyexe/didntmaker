package handlers

import (
	"fmt"
	"log/slog"

	"weazyexe.dev/didntmaker/internal/service"

	tele "gopkg.in/telebot.v3"
)

func (h *Handlers) StatsDay(c tele.Context) error {
	logCommand(c, "/stats_day")
	return h.handleStatsPeriod(c, service.PeriodDay, "день")
}

func (h *Handlers) StatsMonth(c tele.Context) error {
	logCommand(c, "/stats_month")
	return h.handleStatsPeriod(c, service.PeriodMonth, "месяц")
}

func (h *Handlers) StatsYear(c tele.Context) error {
	logCommand(c, "/stats_year")
	return h.handleStatsPeriod(c, service.PeriodYear, "год")
}

func (h *Handlers) handleStatsPeriod(c tele.Context, period service.Period, periodName string) error {
	chatID := c.Chat().ID
	telegramID := c.Sender().ID

	stats, err := h.transactionService.GetUserStats(chatID, telegramID, period)
	if err != nil {
		slog.Error("failed to get user stats",
			"chat_id", chatID,
			"telegram_id", telegramID,
			"period", period,
			"error", err,
		)
		return c.Reply(h.msg.StatsPeriodError)
	}

	if stats.PlusCount == 0 && stats.MinusCount == 0 {
		return c.Reply(h.msg.StatsPeriodEmpty)
	}

	msg := fmt.Sprintf(h.msg.StatsPeriodHeader, periodName)
	msg += fmt.Sprintf(h.msg.StatsPeriodPlusCount, stats.PlusCount)
	msg += fmt.Sprintf(h.msg.StatsPeriodMinusCount, stats.MinusCount)
	msg += fmt.Sprintf(h.msg.StatsPeriodRatio, stats.PlusPercent, stats.MinusPercent)
	msg += fmt.Sprintf(h.msg.StatsPeriodTotalPlus, stats.TotalPlusSum)
	msg += fmt.Sprintf(h.msg.StatsPeriodTotalMinus, stats.TotalMinusSum)

	if len(stats.TopPlusers) > 0 {
		msg += h.msg.StatsPeriodTopPlusers
		for _, sender := range stats.TopPlusers {
			name := sender.Username
			if name == "" {
				name = fmt.Sprintf("id:%d", sender.TelegramID)
			}
			msg += fmt.Sprintf(h.msg.StatsPeriodTopEntry, name, sender.Count, sender.Total)
		}
	}

	if len(stats.TopMinusers) > 0 {
		msg += h.msg.StatsPeriodTopMinusers
		for _, sender := range stats.TopMinusers {
			name := sender.Username
			if name == "" {
				name = fmt.Sprintf("id:%d", sender.TelegramID)
			}
			msg += fmt.Sprintf(h.msg.StatsPeriodTopEntry, name, sender.Count, sender.Total)
		}
	}

	return c.Reply(msg)
}
