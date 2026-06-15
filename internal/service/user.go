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

	now := nowUTC()
	weekDelta, err := s.postingsRepository.ScoreSince(ctx, chatID, telegramID, now.AddDate(0, 0, -7))
	if err != nil {
		return nil, err
	}

	monthDelta, err := s.postingsRepository.ScoreSince(ctx, chatID, telegramID, now.AddDate(0, 0, -30))
	if err != nil {
		return nil, err
	}

	incoming, err := s.postingsRepository.IncomingByCounterparty(ctx, chatID, telegramID)
	if err != nil {
		return nil, err
	}

	outgoing, err := s.postingsRepository.OutgoingByAccount(ctx, chatID, telegramID)
	if err != nil {
		return nil, err
	}

	fan, hater := topPlus(incoming), topMinus(incoming)
	favorite, victim := topPlus(outgoing), topMinus(outgoing)

	return &domain.UserStats{
		User:           user,
		Score:          score,
		DailyRemaining: s.dailyLimit - netSpent,
		DailyLimit:     s.dailyLimit,
		Won:            won,
		Lost:           lost,
		BetAvailable:   !betUsed,
		WeekDelta:      weekDelta,
		MonthDelta:     monthDelta,
		Fan:            fan,
		Hater:          hater,
		Favorite:       favorite,
		Victim:         victim,
	}, nil
}

// topPlus returns the counterparty with the largest positive total, or nil if none.
func topPlus(aggs []domain.CounterpartyAgg) *domain.Counterparty {
	var best *domain.CounterpartyAgg
	for i := range aggs {
		if aggs[i].Plus > 0 && (best == nil || aggs[i].Plus > best.Plus) {
			best = &aggs[i]
		}
	}
	if best == nil {
		return nil
	}
	return &domain.Counterparty{Username: best.Username, FirstName: best.FirstName, Amount: best.Plus}
}

// topMinus returns the counterparty with the largest negative total, or nil if none.
func topMinus(aggs []domain.CounterpartyAgg) *domain.Counterparty {
	var best *domain.CounterpartyAgg
	for i := range aggs {
		if aggs[i].Minus > 0 && (best == nil || aggs[i].Minus > best.Minus) {
			best = &aggs[i]
		}
	}
	if best == nil {
		return nil
	}
	return &domain.Counterparty{Username: best.Username, FirstName: best.FirstName, Amount: best.Minus}
}

func (s *userService) GetLeaderboard(ctx context.Context, chatID int64) ([]domain.LeaderboardEntry, error) {
	return s.postingsRepository.Leaderboard(ctx, chatID)
}

func (s *userService) DailyLimit() int64 {
	return s.dailyLimit
}
