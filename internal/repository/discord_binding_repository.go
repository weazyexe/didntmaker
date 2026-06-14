package repository

import (
	"context"
	"database/sql"
	"time"

	"weazyexe.dev/didntmaker/internal/database/gen"
	"weazyexe.dev/didntmaker/internal/domain"
)

type DiscordBindingRepository interface {
	Create(ctx context.Context, chatID int64, guildID string) error
	Delete(ctx context.Context, chatID int64, guildID string) error
	GetByGuildID(ctx context.Context, guildID string) ([]domain.DiscordBinding, error)
	GetByChatID(ctx context.Context, chatID int64) ([]domain.DiscordBinding, error)
	Exists(ctx context.Context, chatID int64, guildID string) (bool, error)
}

type discordBindingRepository struct {
	queries *gen.Queries
}

func NewDiscordBindingRepository(db *sql.DB) *discordBindingRepository {
	return &discordBindingRepository{queries: gen.New(db)}
}

func (r *discordBindingRepository) Create(ctx context.Context, chatID int64, guildID string) error {
	return r.queries.CreateDiscordBinding(ctx, gen.CreateDiscordBindingParams{
		ChatID:    chatID,
		GuildID:   guildID,
		CreatedAt: time.Now().UTC(),
	})
}

func (r *discordBindingRepository) Delete(ctx context.Context, chatID int64, guildID string) error {
	return r.queries.DeleteDiscordBinding(ctx, gen.DeleteDiscordBindingParams{
		ChatID:  chatID,
		GuildID: guildID,
	})
}

func (r *discordBindingRepository) GetByGuildID(ctx context.Context, guildID string) ([]domain.DiscordBinding, error) {
	rows, err := r.queries.GetDiscordBindingsByGuildID(ctx, guildID)
	if err != nil {
		return nil, err
	}
	return toDomainBindings(rows), nil
}

func (r *discordBindingRepository) GetByChatID(ctx context.Context, chatID int64) ([]domain.DiscordBinding, error) {
	rows, err := r.queries.GetDiscordBindingsByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return toDomainBindings(rows), nil
}

func (r *discordBindingRepository) Exists(ctx context.Context, chatID int64, guildID string) (bool, error) {
	exists, err := r.queries.DiscordBindingExists(ctx, gen.DiscordBindingExistsParams{
		ChatID:  chatID,
		GuildID: guildID,
	})
	if err != nil {
		return false, err
	}
	return exists != 0, nil
}

func toDomainBindings(rows []gen.DiscordBinding) []domain.DiscordBinding {
	bindings := make([]domain.DiscordBinding, 0, len(rows))
	for _, row := range rows {
		bindings = append(bindings, domain.DiscordBinding{
			ChatID:    row.ChatID,
			GuildID:   row.GuildID,
			CreatedAt: row.CreatedAt,
		})
	}
	return bindings
}
