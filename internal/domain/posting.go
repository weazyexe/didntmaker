package domain

import "time"

type Book string

const (
	BookScore     Book = "score"     // очки несправления (/me, /stats)
	BookAllowance Book = "allowance" // дневной лимит раздачи, тратится отправителем
)

type OpType string

const (
	OpTransfer    OpType = "transfer"
	OpTransferAll OpType = "transfer_all"
	OpBetWin      OpType = "bet_win"
	OpBetLose     OpType = "bet_lose"
	OpAdminAdjust OpType = "admin_adjust"
	OpMigration   OpType = "migration"
)

type Posting struct {
	ChatID       int64
	AccountID    int64 // telegram_id of the affected account
	Book         Book
	Amount       int64
	OpID         string
	OpType       OpType
	Counterparty int64 // the other participant, for stats (0 if none)
	Metadata     string
	CreatedAt    time.Time
}
