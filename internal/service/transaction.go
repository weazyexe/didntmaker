package service

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"weazyexe.dev/didntmaker/internal/models"
	"weazyexe.dev/didntmaker/internal/repository"
)

type Period string

const (
	PeriodDay   Period = "day"
	PeriodMonth Period = "month"
	PeriodYear  Period = "year"
)

type TransactionService interface {
	LogTransfer(chatID, senderID, receiverID int64, amount int64)
	LogTransferAll(chatID, senderID int64, amount int64, affectedCount int)
	LogBetResult(chatID, playerID int64, won bool, amount int64, diceValue int)
	LogAdminAdjust(chatID, adminID, targetID int64, delta int64)
	GetUserStats(chatID, telegramID int64, period Period) (*models.UserPeriodStats, error)
	WithTx(tx *gorm.DB) TransactionService
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *transactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) WithTx(tx *gorm.DB) TransactionService {
	return &transactionService{repo: s.repo.WithTx(tx)}
}

func (s *transactionService) LogTransfer(chatID, senderID, receiverID int64, amount int64) {
	tx := &models.Transaction{
		ChatID:     chatID,
		Type:       models.TransactionTypeTransfer,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
	}

	if err := s.repo.Create(tx); err != nil {
		slog.Error("failed to log transfer", "error", err)
	}
}

func (s *transactionService) LogTransferAll(chatID, senderID int64, amount int64, affectedCount int) {
	tx := &models.Transaction{
		ChatID:        chatID,
		Type:          models.TransactionTypeTransferAll,
		SenderID:      senderID,
		ReceiverID:    0,
		Amount:        amount,
		AffectedCount: affectedCount,
	}

	if err := s.repo.Create(tx); err != nil {
		slog.Error("failed to log transfer all", "error", err)
	}
}

func (s *transactionService) LogBetResult(chatID, playerID int64, won bool, amount int64, diceValue int) {
	txType := models.TransactionTypeBetLose
	if won {
		txType = models.TransactionTypeBetWin
	}

	tx := &models.Transaction{
		ChatID:     chatID,
		Type:       txType,
		SenderID:   0,
		ReceiverID: playerID,
		Amount:     amount,
		Metadata:   fmt.Sprintf("dice:%d", diceValue),
	}

	if err := s.repo.Create(tx); err != nil {
		slog.Error("failed to log bet result", "error", err)
	}
}

func (s *transactionService) LogAdminAdjust(chatID, adminID, targetID int64, delta int64) {
	tx := &models.Transaction{
		ChatID:     chatID,
		Type:       models.TransactionTypeAdminAdjust,
		SenderID:   adminID,
		ReceiverID: targetID,
		Amount:     delta,
	}

	if err := s.repo.Create(tx); err != nil {
		slog.Error("failed to log admin adjust", "error", err)
	}
}

func (s *transactionService) GetUserStats(chatID, telegramID int64, period Period) (*models.UserPeriodStats, error) {
	from, to := s.getPeriodBounds(period)

	transactions, err := s.repo.GetByReceiverInPeriod(chatID, telegramID, from, to)
	if err != nil {
		return nil, err
	}

	stats := &models.UserPeriodStats{}

	for _, tx := range transactions {
		if tx.Amount > 0 {
			stats.PlusCount++
			stats.TotalPlusSum += tx.Amount
		} else if tx.Amount < 0 {
			stats.MinusCount++
			stats.TotalMinusSum += tx.Amount
		}
	}

	total := stats.PlusCount + stats.MinusCount
	if total > 0 {
		stats.PlusPercent = float64(stats.PlusCount) / float64(total) * 100
		stats.MinusPercent = float64(stats.MinusCount) / float64(total) * 100
	}

	topPlusers, err := s.repo.GetTopSendersToUser(chatID, telegramID, from, to, true, 3)
	if err != nil {
		slog.Error("failed to get top plusers", "error", err)
	} else {
		stats.TopPlusers = topPlusers
	}

	topMinusers, err := s.repo.GetTopSendersToUser(chatID, telegramID, from, to, false, 3)
	if err != nil {
		slog.Error("failed to get top minusers", "error", err)
	} else {
		stats.TopMinusers = topMinusers
	}

	return stats, nil
}

func (s *transactionService) getPeriodBounds(period Period) (from, to time.Time) {
	now := time.Now().UTC()

	switch period {
	case PeriodDay:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 0, 1)
	case PeriodMonth:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 1, 0)
	case PeriodYear:
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(1, 0, 0)
	default:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 0, 1)
	}

	return from, to
}
