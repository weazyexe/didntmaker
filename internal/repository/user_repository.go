package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"weazyexe.dev/didntmaker/internal/database/gen"
	"weazyexe.dev/didntmaker/internal/domain"
)

type UserRepository interface {
	GetOrCreate(ctx context.Context, chatID, telegramID int64, username, firstName string) (domain.User, error)
	GetByUsername(ctx context.Context, chatID int64, username string) (domain.User, error)
	GetByTelegramID(ctx context.Context, chatID, telegramID int64) (domain.User, error)
	ListChatUsers(ctx context.Context, chatID int64) ([]domain.User, error)
	CountUsersExcept(ctx context.Context, chatID, exceptTelegramID int64) (int64, error)
	ListUserIDsExcept(ctx context.Context, chatID, exceptTelegramID int64) ([]int64, error)
}

type userRepository struct {
	queries *gen.Queries
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{queries: gen.New(db)}
}

func (r *userRepository) GetOrCreate(ctx context.Context, chatID, telegramID int64, username, firstName string) (domain.User, error) {
	row, err := r.queries.GetUserByTelegramID(ctx, gen.GetUserByTelegramIDParams{
		ChatID:     chatID,
		TelegramID: telegramID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		now := time.Now().UTC()
		created, err := r.queries.CreateUser(ctx, gen.CreateUserParams{
			TelegramID: telegramID,
			ChatID:     chatID,
			Username:   username,
			FirstName:  firstName,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		if err != nil {
			return domain.User{}, err
		}
		return toDomainUser(created), nil
	}
	if err != nil {
		return domain.User{}, err
	}

	// Refresh profile fields if they changed.
	if row.Username != username || row.FirstName != firstName {
		if err := r.queries.UpdateUserProfile(ctx, gen.UpdateUserProfileParams{
			Username:   username,
			FirstName:  firstName,
			UpdatedAt:  time.Now().UTC(),
			ChatID:     chatID,
			TelegramID: telegramID,
		}); err != nil {
			return domain.User{}, err
		}
		row.Username = username
		row.FirstName = firstName
	}

	return toDomainUser(row), nil
}

func (r *userRepository) GetByUsername(ctx context.Context, chatID int64, username string) (domain.User, error) {
	row, err := r.queries.GetUserByUsername(ctx, gen.GetUserByUsernameParams{ChatID: chatID, Username: username})
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(row), nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, chatID, telegramID int64) (domain.User, error) {
	row, err := r.queries.GetUserByTelegramID(ctx, gen.GetUserByTelegramIDParams{ChatID: chatID, TelegramID: telegramID})
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(row), nil
}

func (r *userRepository) ListChatUsers(ctx context.Context, chatID int64) ([]domain.User, error) {
	rows, err := r.queries.ListChatUsers(ctx, chatID)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toDomainUser(row))
	}
	return users, nil
}

func (r *userRepository) CountUsersExcept(ctx context.Context, chatID, exceptTelegramID int64) (int64, error) {
	return r.queries.CountChatUsersExcept(ctx, gen.CountChatUsersExceptParams{
		ChatID:     chatID,
		TelegramID: exceptTelegramID,
	})
}

func (r *userRepository) ListUserIDsExcept(ctx context.Context, chatID, exceptTelegramID int64) ([]int64, error) {
	return r.queries.ListChatUserIDsExcept(ctx, gen.ListChatUserIDsExceptParams{
		ChatID:     chatID,
		TelegramID: exceptTelegramID,
	})
}

func toDomainUser(u gen.User) domain.User {
	return domain.User{
		ID:         u.ID,
		TelegramID: u.TelegramID,
		ChatID:     u.ChatID,
		Username:   u.Username,
		FirstName:  u.FirstName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
