package repository

import (
	"log/slog"
	"time"

	"weazyexe.dev/didntmaker/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db         *gorm.DB
	dailyLimit int64
}

func NewUserRepository(db *gorm.DB, dailyLimit int64) *UserRepository {
	slog.Info("user repository created", "daily_limit", dailyLimit)
	return &UserRepository{db: db, dailyLimit: dailyLimit}
}

func (r *UserRepository) DailyLimit() int64 {
	return r.dailyLimit
}

func (r *UserRepository) GetOrCreateUser(chatID, telegramID int64, username, firstName string) (*models.User, error) {
	var user models.User
	result := r.db.Where("telegram_id = ? AND chat_id = ?", telegramID, chatID).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		user = models.User{
			TelegramID: telegramID,
			ChatID:     chatID,
			Username:   username,
			FirstName:  firstName,
			Balance:    0,
		}
		if err := r.db.Create(&user).Error; err != nil {
			slog.Error("failed to create user", "chat_id", chatID, "telegram_id", telegramID, "error", err)
			return nil, err
		}
		slog.Info("user created", "chat_id", chatID, "telegram_id", telegramID, "username", username)
		return &user, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	if user.Username != username || user.FirstName != firstName {
		user.Username = username
		user.FirstName = firstName
		r.db.Save(&user)
		slog.Debug("user info updated", "chat_id", chatID, "telegram_id", telegramID, "username", username)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(chatID int64, username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("chat_id = ? AND username = ?", chatID, username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByTelegramID(chatID, telegramID int64) (*models.User, error) {
	var user models.User
	result := r.db.Where("chat_id = ? AND telegram_id = ?", chatID, telegramID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) UpdateBalance(chatID, telegramID int64, delta int64) (*models.User, int64, error) {
	var user models.User
	result := r.db.Where("telegram_id = ? AND chat_id = ?", telegramID, chatID).First(&user)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	oldBalance := user.Balance
	user.Balance += delta

	if err := r.db.Save(&user).Error; err != nil {
		slog.Error("failed to update balance", "chat_id", chatID, "telegram_id", telegramID, "error", err)
		return nil, 0, err
	}

	slog.Info("balance updated",
		"chat_id", chatID,
		"telegram_id", telegramID,
		"old_balance", oldBalance,
		"new_balance", user.Balance,
		"delta", delta,
	)

	return &user, oldBalance, nil
}

func (r *UserRepository) UpdateBalanceByUsername(chatID int64, username string, delta int64) (*models.User, int64, error) {
	var user models.User
	result := r.db.Where("chat_id = ? AND username = ?", chatID, username).First(&user)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	oldBalance := user.Balance
	user.Balance += delta

	if err := r.db.Save(&user).Error; err != nil {
		slog.Error("failed to update balance by username", "chat_id", chatID, "username", username, "error", err)
		return nil, 0, err
	}

	slog.Info("balance updated by username",
		"chat_id", chatID,
		"username", username,
		"old_balance", oldBalance,
		"new_balance", user.Balance,
		"delta", delta,
	)

	return &user, oldBalance, nil
}

func (r *UserRepository) GetChatStats(chatID int64) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("chat_id = ?", chatID).Order("balance ASC").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepository) UpdateBalanceForAllExcept(chatID, exceptTelegramID int64, delta int64) (int, error) {
	result := r.db.Model(&models.User{}).
		Where("chat_id = ? AND telegram_id != ?", chatID, exceptTelegramID).
		Update("balance", gorm.Expr("balance + ?", delta))
	if result.Error != nil {
		slog.Error("failed to update balance for all", "chat_id", chatID, "error", result.Error)
		return 0, result.Error
	}

	slog.Info("balance updated for all",
		"chat_id", chatID,
		"except_telegram_id", exceptTelegramID,
		"delta", delta,
		"affected", result.RowsAffected,
	)

	return int(result.RowsAffected), nil
}

func (r *UserRepository) CountUsersExcept(chatID, exceptTelegramID int64) (int64, error) {
	var count int64
	result := r.db.Model(&models.User{}).
		Where("chat_id = ? AND telegram_id != ?", chatID, exceptTelegramID).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *UserRepository) GetDailyRemaining(chatID, telegramID int64) (int64, error) {
	var user models.User
	result := r.db.Where("telegram_id = ? AND chat_id = ?", telegramID, chatID).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	resetDay := user.DailyResetAt.Truncate(24 * time.Hour)

	if today.After(resetDay) {
		user.DailyGiven = 0
		user.DailyResetAt = time.Now().UTC()
		r.db.Save(&user)
		slog.Info("daily limit reset", "chat_id", chatID, "telegram_id", telegramID)
		return r.dailyLimit, nil
	}

	return r.dailyLimit - user.DailyGiven, nil
}

func (r *UserRepository) AddDailyGiven(chatID, telegramID int64, amount int64) error {
	var user models.User
	result := r.db.Where("telegram_id = ? AND chat_id = ?", telegramID, chatID).First(&user)
	if result.Error != nil {
		return result.Error
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	resetDay := user.DailyResetAt.Truncate(24 * time.Hour)

	if today.After(resetDay) {
		user.DailyGiven = amount
		user.DailyResetAt = time.Now().UTC()
	} else {
		user.DailyGiven += amount
	}

	slog.Debug("daily given updated",
		"chat_id", chatID,
		"telegram_id", telegramID,
		"amount", amount,
		"total_given", user.DailyGiven,
	)

	return r.db.Save(&user).Error
}

func (r *UserRepository) CanBetToday(chatID, telegramID int64) (bool, error) {
	var user models.User
	result := r.db.Where("telegram_id = ? AND chat_id = ?", telegramID, chatID).First(&user)
	if result.Error != nil {
		return false, result.Error
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastBetDay := user.LastBetAt.Truncate(24 * time.Hour)

	return today.After(lastBetDay), nil
}

func (r *UserRepository) ApplyBetResult(chatID, telegramID int64, won bool) error {
	var user models.User
	result := r.db.Where("telegram_id = ? AND chat_id = ?", telegramID, chatID).First(&user)
	if result.Error != nil {
		return result.Error
	}

	user.LastBetAt = time.Now().UTC()

	if won {
		user.DailyGiven -= r.dailyLimit
		if user.DailyGiven < 0 {
			user.DailyGiven = 0
		}
		slog.Info("bet won", "chat_id", chatID, "telegram_id", telegramID, "daily_limit_bonus", r.dailyLimit)
	} else {
		user.Balance -= r.dailyLimit
		slog.Info("bet lost", "chat_id", chatID, "telegram_id", telegramID, "balance_penalty", r.dailyLimit, "new_balance", user.Balance)
	}

	return r.db.Save(&user).Error
}

func (r *UserRepository) AddDailyLimitByUsername(chatID int64, username string, delta int64) (oldRemaining, newRemaining int64, err error) {
	var user models.User
	result := r.db.Where("chat_id = ? AND username = ?", chatID, username).First(&user)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	resetDay := user.DailyResetAt.Truncate(24 * time.Hour)

	if today.After(resetDay) {
		user.DailyGiven = 0
		user.DailyResetAt = time.Now().UTC()
	}

	oldRemaining = r.dailyLimit - user.DailyGiven

	user.DailyGiven -= delta
	if user.DailyGiven < 0 {
		user.DailyGiven = 0
	}

	newRemaining = r.dailyLimit - user.DailyGiven

	slog.Info("admin adjusted daily limit",
		"chat_id", chatID,
		"username", username,
		"delta", delta,
		"old_remaining", oldRemaining,
		"new_remaining", newRemaining,
	)

	err = r.db.Save(&user).Error
	return oldRemaining, newRemaining, err
}
