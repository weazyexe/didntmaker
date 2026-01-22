package service

import (
	"strings"
	"time"

	"weazyexe.dev/didntmaker/internal/models"
	"weazyexe.dev/didntmaker/internal/repository"
)

const maxTransactionAmount int64 = 1000

type TransferResult struct {
	Target     *models.User
	Delta      int64
	OldBalance int64
}

type TransferAllResult struct {
	Delta       int64
	AffectedCnt int
	TotalCost   int64
}

type DailyBalance struct {
	User       *models.User
	Remaining  int64
	DailyLimit int64
}

type AdjustResult struct {
	OldRemaining int64
	NewRemaining int64
}

type BalanceService struct {
	repo       *repository.UserRepository
	superAdmin string
}

func NewBalanceService(repo *repository.UserRepository, superAdmin string) *BalanceService {
	return &BalanceService{repo: repo, superAdmin: superAdmin}
}

func (s *BalanceService) Transfer(chatID, senderID, targetID int64, targetUsername string, delta int64) (*TransferResult, error) {
	if senderID == targetID {
		return nil, ErrSelfTransfer
	}

	if delta > maxTransactionAmount || delta < -maxTransactionAmount {
		return nil, ErrTransactionLimit
	}

	absDelta := abs(delta)

	remaining, err := s.repo.GetDailyRemaining(chatID, senderID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if absDelta > remaining {
		return nil, ErrInsufficientLimit
	}

	user, oldBalance, err := s.repo.UpdateBalance(chatID, targetID, delta)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") && targetUsername != "" {
			user, oldBalance, err = s.repo.UpdateBalanceByUsername(chatID, targetUsername, delta)
		}
		if err != nil {
			return nil, ErrUserNotFound
		}
	}

	_ = s.repo.AddDailyGiven(chatID, senderID, absDelta)

	return &TransferResult{
		Target:     user,
		Delta:      delta,
		OldBalance: oldBalance,
	}, nil
}

func (s *BalanceService) TransferToAll(chatID, senderID int64, delta int64) (*TransferAllResult, error) {
	if delta > maxTransactionAmount || delta < -maxTransactionAmount {
		return nil, ErrTransactionLimit
	}

	userCount, err := s.repo.CountUsersExcept(chatID, senderID)
	if err != nil || userCount == 0 {
		return nil, ErrNoUsersInChat
	}

	absDelta := abs(delta)
	totalCost := absDelta * userCount

	remaining, err := s.repo.GetDailyRemaining(chatID, senderID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if totalCost > remaining {
		return nil, ErrInsufficientLimit
	}

	affected, err := s.repo.UpdateBalanceForAllExcept(chatID, senderID, delta)
	if err != nil {
		return nil, err
	}

	_ = s.repo.AddDailyGiven(chatID, senderID, totalCost)

	return &TransferAllResult{
		Delta:       delta,
		AffectedCnt: affected,
		TotalCost:   totalCost,
	}, nil
}

func (s *BalanceService) GetDailyBalances(chatID int64) ([]DailyBalance, error) {
	users, err := s.repo.GetChatStats(chatID)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	dailyLimit := s.repo.DailyLimit()

	result := make([]DailyBalance, 0, len(users))
	for _, user := range users {
		var remaining int64
		resetDay := user.DailyResetAt.Truncate(24 * time.Hour)
		if today.After(resetDay) {
			remaining = dailyLimit
		} else {
			remaining = dailyLimit - user.DailyGiven
		}

		u := user // copy to avoid pointer issues
		result = append(result, DailyBalance{
			User:       &u,
			Remaining:  remaining,
			DailyLimit: dailyLimit,
		})
	}

	return result, nil
}

func (s *BalanceService) AdjustDailyLimit(chatID int64, adminUsername, targetUsername string, delta int64) (*AdjustResult, error) {
	if s.superAdmin == "" || !strings.EqualFold(adminUsername, s.superAdmin) {
		return nil, ErrNotAuthorized
	}

	oldRemaining, newRemaining, err := s.repo.AddDailyLimitByUsername(chatID, targetUsername, delta)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &AdjustResult{
		OldRemaining: oldRemaining,
		NewRemaining: newRemaining,
	}, nil
}

func (s *BalanceService) DailyLimit() int64 {
	return s.repo.DailyLimit()
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
