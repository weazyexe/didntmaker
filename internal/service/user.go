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

type UserService interface {
	GetOrCreate(chatID, telegramID int64, username, firstName string) (*models.User, error)
	GetStats(chatID, telegramID int64, username, firstName string) (*UserStats, error)
	GetLeaderboard(chatID int64) ([]models.User, error)
	DailyLimit() int64
}

// userService implements UserSvc interface
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetOrCreate(chatID, telegramID int64, username, firstName string) (*models.User, error) {
	return s.repo.GetOrCreateUser(chatID, telegramID, username, firstName)
}

func (s *userService) GetStats(chatID, telegramID int64, username, firstName string) (*UserStats, error) {
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

func (s *userService) GetLeaderboard(chatID int64) ([]models.User, error) {
	return s.repo.GetChatStats(chatID)
}

func (s *userService) DailyLimit() int64 {
	return s.repo.DailyLimit()
}
