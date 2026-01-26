package models

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeTransfer    TransactionType = "transfer"
	TransactionTypeTransferAll TransactionType = "transfer_all"
	TransactionTypeBetWin      TransactionType = "bet_win"
	TransactionTypeBetLose     TransactionType = "bet_lose"
	TransactionTypeAdminAdjust TransactionType = "admin_adjust"
	TransactionTypeMigration   TransactionType = "migration"
)

type Transaction struct {
	ID            uint            `gorm:"primaryKey"`
	ChatID        int64           `gorm:"index:idx_chat_created"`
	Type          TransactionType `gorm:"type:varchar(20);index"`
	SenderID      int64           `gorm:"index"`
	ReceiverID    int64           `gorm:"index"`
	Amount        int64
	AffectedCount int
	Metadata      string
	CreatedAt     time.Time `gorm:"index:idx_chat_created"`
}

type SenderStats struct {
	TelegramID int64
	Username   string
	Count      int64
	Total      int64
}

type UserPeriodStats struct {
	PlusCount     int64
	MinusCount    int64
	PlusPercent   float64
	MinusPercent  float64
	TotalPlusSum  int64
	TotalMinusSum int64
	TopPlusers    []SenderStats
	TopMinusers   []SenderStats
}
