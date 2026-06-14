package service

import (
	"context"

	"weazyexe.dev/didntmaker/internal/domain"
	"weazyexe.dev/didntmaker/internal/repository"
)

type UserService interface {
	GetOrCreate(ctx context.Context, chatID, telegramID int64, username, firstName string) (domain.User, error)
	GetStats(ctx context.Context, chatID, telegramID int64, username, firstName string) (*domain.UserStats, error)
	GetLeaderboard(ctx context.Context, chatID int64) ([]domain.LeaderboardEntry, error)
	DailyLimit() int64
}

type userService struct {
	usersRepository    repository.UserRepository
	postingsRepository repository.PostingRepository
	dailyLimit         int64
}

func NewUserService(users repository.UserRepository, postings repository.PostingRepository, dailyLimit int64) *userService {
	return &userService{usersRepository: users, postingsRepository: postings, dailyLimit: dailyLimit}
}

func (s *userService) GetOrCreate(ctx context.Context, chatID, telegramID int64, username, firstName string) (domain.User, error) {
	return s.usersRepository.GetOrCreate(ctx, chatID, telegramID, username, firstName)
}

func (s *userService) GetStats(ctx context.Context, chatID, telegramID int64, username, firstName string) (*domain.UserStats, error) {
	user, err := s.usersRepository.GetOrCreate(ctx, chatID, telegramID, username, firstName)
	if err != nil {
		return nil, err
	}

	score, err := s.postingsRepository.Score(ctx, chatID, telegramID)
	if err != nil {
		return nil, err
	}

	netSpent, err := s.postingsRepository.AllowanceSpentSince(ctx, chatID, telegramID, startOfUTCDay())
	if err != nil {
		return nil, err
	}

	won, lost, err := s.postingsRepository.BetStats(ctx, chatID, telegramID)
	if err != nil {
		return nil, err
	}

	betUsed, err := s.postingsRepository.HasBetSince(ctx, chatID, telegramID, startOfUTCDay())
	if err != nil {
		return nil, err
	}

	return &domain.UserStats{
		User:           user,
		Score:          score,
		DailyRemaining: s.dailyLimit - netSpent,
		DailyLimit:     s.dailyLimit,
		Won:            won,
		Lost:           lost,
		BetAvailable:   !betUsed,
	}, nil
}

func (s *userService) GetLeaderboard(ctx context.Context, chatID int64) ([]domain.LeaderboardEntry, error) {
	return s.postingsRepository.Leaderboard(ctx, chatID)
}

func (s *userService) DailyLimit() int64 {
	return s.dailyLimit
}
