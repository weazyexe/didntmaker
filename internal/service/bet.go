package service

import (
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

// betService implements BetSvc interface
type betService struct {
	repo repository.UserRepository
}

func NewBetService(repo repository.UserRepository) *betService {
	return &betService{repo: repo}
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

	if err := s.repo.ApplyBetResult(chatID, telegramID, won); err != nil {
		return nil, err
	}

	return &BetResult{
		DiceValue:  diceValue,
		Won:        won,
		DailyLimit: s.repo.DailyLimit(),
	}, nil
}

func (s *betService) DailyLimit() int64 {
	return s.repo.DailyLimit()
}
