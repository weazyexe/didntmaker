package repository

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
	"weazyexe.dev/didntmaker/internal/models"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	GetByReceiverInPeriod(chatID, receiverID int64, from, to time.Time) ([]models.Transaction, error)
	GetTopSendersToUser(chatID, receiverID int64, from, to time.Time, positive bool, limit int) ([]models.SenderStats, error)
	WithTx(tx *gorm.DB) TransactionRepository
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	slog.Info("transaction repository created")
	return &transactionRepository{db: db}
}

func (r *transactionRepository) WithTx(tx *gorm.DB) TransactionRepository {
	return &transactionRepository{db: tx}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = time.Now().UTC()
	}

	if err := r.db.Create(tx).Error; err != nil {
		slog.Error("failed to create transaction",
			"chat_id", tx.ChatID,
			"type", tx.Type,
			"sender_id", tx.SenderID,
			"receiver_id", tx.ReceiverID,
			"amount", tx.Amount,
			"error", err,
		)
		return err
	}

	slog.Debug("transaction created",
		"id", tx.ID,
		"chat_id", tx.ChatID,
		"type", tx.Type,
		"sender_id", tx.SenderID,
		"receiver_id", tx.ReceiverID,
		"amount", tx.Amount,
	)

	return nil
}

func (r *transactionRepository) GetByReceiverInPeriod(chatID, receiverID int64, from, to time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Where("chat_id = ? AND receiver_id = ? AND created_at >= ? AND created_at < ?",
		chatID, receiverID, from, to).
		Order("created_at DESC").
		Find(&transactions)

	if result.Error != nil {
		slog.Error("failed to get transactions by receiver",
			"chat_id", chatID,
			"receiver_id", receiverID,
			"from", from,
			"to", to,
			"error", result.Error,
		)
		return nil, result.Error
	}

	return transactions, nil
}

func (r *transactionRepository) GetTopSendersToUser(chatID, receiverID int64, from, to time.Time, positive bool, limit int) ([]models.SenderStats, error) {
	var stats []models.SenderStats

	amountCondition := "amount > 0"
	orderDirection := "DESC"
	if !positive {
		amountCondition = "amount < 0"
		orderDirection = "ASC"
	}

	query := `
		SELECT
			t.sender_id as telegram_id,
			COALESCE(u.username, '') as username,
			COUNT(*) as count,
			SUM(t.amount) as total
		FROM transactions t
		LEFT JOIN users u ON t.chat_id = u.chat_id AND t.sender_id = u.telegram_id
		WHERE t.chat_id = ?
			AND t.receiver_id = ?
			AND t.created_at >= ?
			AND t.created_at < ?
			AND t.sender_id != 0
			AND ` + amountCondition + `
		GROUP BY t.sender_id, u.username
		ORDER BY total ` + orderDirection + `
		LIMIT ?
	`

	result := r.db.Raw(query, chatID, receiverID, from, to, limit).Scan(&stats)
	if result.Error != nil {
		slog.Error("failed to get top senders",
			"chat_id", chatID,
			"receiver_id", receiverID,
			"positive", positive,
			"error", result.Error,
		)
		return nil, result.Error
	}

	return stats, nil
}
