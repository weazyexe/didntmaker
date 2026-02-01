package repository

import (
	"log/slog"

	"gorm.io/gorm"
	"weazyexe.dev/didntmaker/internal/models"
)

type DiscordBindingRepository interface {
	Create(binding *models.DiscordBinding) error
	Delete(chatID int64, guildID string) error
	GetByGuildID(guildID string) ([]models.DiscordBinding, error)
	GetByChatID(chatID int64) ([]models.DiscordBinding, error)
	Exists(chatID int64, guildID string) (bool, error)
}

type discordBindingRepository struct {
	db *gorm.DB
}

func NewDiscordBindingRepository(db *gorm.DB) *discordBindingRepository {
	slog.Info("discord binding repository created")
	return &discordBindingRepository{db: db}
}

func (r *discordBindingRepository) Create(binding *models.DiscordBinding) error {
	if err := r.db.Create(binding).Error; err != nil {
		slog.Error("failed to create discord binding",
			"chat_id", binding.ChatID,
			"guild_id", binding.GuildID,
			"error", err,
		)
		return err
	}

	slog.Info("discord binding created",
		"chat_id", binding.ChatID,
		"guild_id", binding.GuildID,
	)
	return nil
}

func (r *discordBindingRepository) Delete(chatID int64, guildID string) error {
	result := r.db.Where("chat_id = ? AND guild_id = ?", chatID, guildID).Delete(&models.DiscordBinding{})
	if result.Error != nil {
		slog.Error("failed to delete discord binding",
			"chat_id", chatID,
			"guild_id", guildID,
			"error", result.Error,
		)
		return result.Error
	}

	slog.Info("discord binding deleted",
		"chat_id", chatID,
		"guild_id", guildID,
		"rows_affected", result.RowsAffected,
	)
	return nil
}

func (r *discordBindingRepository) GetByGuildID(guildID string) ([]models.DiscordBinding, error) {
	var bindings []models.DiscordBinding
	if err := r.db.Where("guild_id = ?", guildID).Find(&bindings).Error; err != nil {
		slog.Error("failed to get discord bindings by guild_id",
			"guild_id", guildID,
			"error", err,
		)
		return nil, err
	}
	return bindings, nil
}

func (r *discordBindingRepository) GetByChatID(chatID int64) ([]models.DiscordBinding, error) {
	var bindings []models.DiscordBinding
	if err := r.db.Where("chat_id = ?", chatID).Find(&bindings).Error; err != nil {
		slog.Error("failed to get discord bindings by chat_id",
			"chat_id", chatID,
			"error", err,
		)
		return nil, err
	}
	return bindings, nil
}

func (r *discordBindingRepository) Exists(chatID int64, guildID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.DiscordBinding{}).
		Where("chat_id = ? AND guild_id = ?", chatID, guildID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
