package service

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"weazyexe.dev/didntmaker/internal/models"
	"weazyexe.dev/didntmaker/internal/repository"
)

const maxTransactionAmount int64 = 1000

type BalanceService interface {
	Transfer(chatID, senderID, targetID int64, targetUsername string, delta int64) (*models.TransferResult, error)
	TransferToAll(chatID, senderID int64, delta int64) (*models.TransferAllResult, error)
	GetDailyBalances(chatID int64) ([]models.DailyBalance, error)
	AdjustDailyLimit(chatID int64, adminUsername, targetUsername string, delta int64) (*models.AdjustResult, error)
	DailyLimit() int64
}

// balanceService implements BalanceSvc interface
type balanceService struct {
	repo       repository.UserRepository
	txService  TransactionService
	superAdmin []string
}

func NewBalanceService(repo repository.UserRepository, txService TransactionService, superAdmin []string) *balanceService {
	return &balanceService{repo: repo, txService: txService, superAdmin: superAdmin}
}

func (s *balanceService) Transfer(chatID, senderID, targetID int64, targetUsername string, delta int64) (*models.TransferResult, error) {
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

	var result *models.TransferResult
	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txSvc := s.txService.WithTx(tx)

		user, oldBalance, err := txRepo.UpdateBalance(chatID, targetID, delta)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") && targetUsername != "" {
				user, oldBalance, err = txRepo.UpdateBalanceByUsername(chatID, targetUsername, delta)
			}
			if err != nil {
				return ErrUserNotFound
			}
		}

		if err := txRepo.AddDailyGiven(chatID, senderID, absDelta); err != nil {
			return err
		}

		txSvc.LogTransfer(chatID, senderID, user.TelegramID, delta)

		result = &models.TransferResult{
			Target:     user,
			Delta:      delta,
			OldBalance: oldBalance,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *balanceService) TransferToAll(chatID, senderID int64, delta int64) (*models.TransferAllResult, error) {
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

	var result *models.TransferAllResult
	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txSvc := s.txService.WithTx(tx)

		affected, err := txRepo.UpdateBalanceForAllExcept(chatID, senderID, delta)
		if err != nil {
			return err
		}

		if err := txRepo.AddDailyGiven(chatID, senderID, totalCost); err != nil {
			return err
		}

		txSvc.LogTransferAll(chatID, senderID, delta, affected)

		result = &models.TransferAllResult{
			Delta:       delta,
			AffectedCnt: affected,
			TotalCost:   totalCost,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *balanceService) GetDailyBalances(chatID int64) ([]models.DailyBalance, error) {
	users, err := s.repo.GetChatStats(chatID)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	dailyLimit := s.repo.DailyLimit()

	result := make([]models.DailyBalance, 0, len(users))
	for _, user := range users {
		var remaining int64
		resetDay := user.DailyResetAt.Truncate(24 * time.Hour)
		if today.After(resetDay) {
			remaining = dailyLimit
		} else {
			remaining = dailyLimit - user.DailyGiven
		}

		u := user // copy to avoid pointer issues
		result = append(result, models.DailyBalance{
			User:       &u,
			Remaining:  remaining,
			DailyLimit: dailyLimit,
		})
	}

	return result, nil
}

func (s *balanceService) AdjustDailyLimit(chatID int64, adminUsername, targetUsername string, delta int64) (*models.AdjustResult, error) {
	if !s.isSuperAdmin(adminUsername) {
		return nil, ErrNotAuthorized
	}

	oldRemaining, newRemaining, err := s.repo.AddDailyLimitByUsername(chatID, targetUsername, delta)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &models.AdjustResult{
		OldRemaining: oldRemaining,
		NewRemaining: newRemaining,
	}, nil
}

func (s *balanceService) isSuperAdmin(username string) bool {
	for _, admin := range s.superAdmin {
		if strings.EqualFold(username, admin) {
			return true
		}
	}
	return false
}

func (s *balanceService) DailyLimit() int64 {
	return s.repo.DailyLimit()
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
