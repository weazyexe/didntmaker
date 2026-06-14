package service

import (
	"context"
	"fmt"
	"log/slog"

	"weazyexe.dev/didntmaker/internal/domain"
	"weazyexe.dev/didntmaker/internal/repository"
)

type BetService interface {
	CanBet(ctx context.Context, chatID, telegramID int64) error
	ApplyResult(ctx context.Context, chatID, telegramID int64, diceValue int) (*domain.BetResult, error)
	DailyLimit() int64
}

type betService struct {
	postingRepository repository.PostingRepository
	dailyLimit        int64
}

func NewBetService(postingRepository repository.PostingRepository, dailyLimit int64) *betService {
	return &betService{postingRepository: postingRepository, dailyLimit: dailyLimit}
}

func (s *betService) CanBet(ctx context.Context, chatID, telegramID int64) error {
	used, err := s.postingRepository.HasBetSince(ctx, chatID, telegramID, startOfUTCDay())
	if err != nil {
		return err
	}
	if used {
		return domain.ErrBetAlreadyUsed
	}
	return nil
}

func (s *betService) ApplyResult(ctx context.Context, chatID, telegramID int64, diceValue int) (*domain.BetResult, error) {
	won := diceValue >= 4

	// Win adds +limit to daily allowance (negative spend, stacks); loss deducts score.
	posting := domain.Posting{
		ChatID:       chatID,
		AccountID:    telegramID,
		Amount:       -s.dailyLimit,
		OpID:         newOpID(),
		Counterparty: 0,
		Metadata:     fmt.Sprintf("dice:%d", diceValue),
		CreatedAt:    nowUTC(),
	}
	if won {
		posting.Book = domain.BookAllowance
		posting.OpType = domain.OpBetWin
	} else {
		posting.Book = domain.BookScore
		posting.OpType = domain.OpBetLose
	}

	if err := s.postingRepository.InsertPostings(ctx, []domain.Posting{posting}); err != nil {
		slog.Error("bet failed", "chat_id", chatID, "player", telegramID, "error", err)
		return nil, err
	}

	slog.Info("bet", "chat_id", chatID, "player", telegramID, "won", won, "dice", diceValue)
	return &domain.BetResult{
		DiceValue:  diceValue,
		Won:        won,
		DailyLimit: s.dailyLimit,
	}, nil
}

func (s *betService) DailyLimit() int64 {
	return s.dailyLimit
}
