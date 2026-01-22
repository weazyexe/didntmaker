package service

import (
	"weazyexe.dev/didntmaker/internal/models"
	"weazyexe.dev/didntmaker/internal/repository"
)

type UserStats struct {
	User           *models.User
	DailyRemaining int64
	DailyLimit     int64
}

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetOrCreate(chatID, telegramID int64, username, firstName string) (*models.User, error) {
	return s.repo.GetOrCreateUser(chatID, telegramID, username, firstName)
}

func (s *UserService) GetStats(chatID, telegramID int64, username, firstName string) (*UserStats, error) {
	user, err := s.repo.GetOrCreateUser(chatID, telegramID, username, firstName)
	if err != nil {
		return nil, err
	}

	remaining, err := s.repo.GetDailyRemaining(chatID, telegramID)
	if err != nil {
		return nil, err
	}

	return &UserStats{
		User:           user,
		DailyRemaining: remaining,
		DailyLimit:     s.repo.DailyLimit(),
	}, nil
}

func (s *UserService) GetLeaderboard(chatID int64) ([]models.User, error) {
	return s.repo.GetChatStats(chatID)
}

func (s *UserService) DailyLimit() int64 {
	return s.repo.DailyLimit()
}
