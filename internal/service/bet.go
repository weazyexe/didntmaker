package service

import (
	"gorm.io/gorm"
	"weazyexe.dev/didntmaker/internal/repository"
)

type BetResult struct {
	DiceValue  int
	Won        bool
	DailyLimit int64
}

type BetService interface {
	CanBet(chatID, telegramID int64) error
	ApplyResult(chatID, telegramID int64, diceValue int) (*BetResult, error)
	DailyLimit() int64
}

type betService struct {
	repo      repository.UserRepository
	txService TransactionService
}

func NewBetService(repo repository.UserRepository, txService TransactionService) *betService {
	return &betService{repo: repo, txService: txService}
}

func (s *betService) CanBet(chatID, telegramID int64) error {
	canBet, err := s.repo.CanBetToday(chatID, telegramID)
	if err != nil {
		return err
	}
	if !canBet {
		return ErrBetAlreadyUsed
	}
	return nil
}

func (s *betService) ApplyResult(chatID, telegramID int64, diceValue int) (*BetResult, error) {
	won := diceValue >= 4
	dailyLimit := s.repo.DailyLimit()

	amount := -dailyLimit
	if won {
		amount = dailyLimit
	}

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txSvc := s.txService.WithTx(tx)

		if err := txRepo.ApplyBetResult(chatID, telegramID, won); err != nil {
			return err
		}

		txSvc.LogBetResult(chatID, telegramID, won, amount, diceValue)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BetResult{
		DiceValue:  diceValue,
		Won:        won,
		DailyLimit: dailyLimit,
	}, nil
}

func (s *betService) DailyLimit() int64 {
	return s.repo.DailyLimit()
}
