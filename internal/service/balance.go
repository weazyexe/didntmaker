package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"weazyexe.dev/didntmaker/internal/domain"
	"weazyexe.dev/didntmaker/internal/repository"
)

type BalanceService interface {
	Transfer(ctx context.Context, chatID, senderID, targetID int64, targetUsername string, delta int64) (*domain.TransferResult, error)
	TransferToAll(ctx context.Context, chatID, senderID int64, delta int64) (*domain.TransferAllResult, error)
	GetDailyBalances(ctx context.Context, chatID int64) ([]domain.DailyBalance, error)
	AdjustDailyLimit(ctx context.Context, chatID int64, adminUsername, targetUsername string, delta int64) (*domain.AdjustResult, error)
	DailyLimit() int64
}

type balanceService struct {
	users      repository.UserRepository
	ledger     repository.PostingRepository
	dailyLimit int64
	superAdmin []string
}

func NewBalanceService(users repository.UserRepository, ledger repository.PostingRepository, dailyLimit int64, superAdmin []string) *balanceService {
	return &balanceService{users: users, ledger: ledger, dailyLimit: dailyLimit, superAdmin: superAdmin}
}

func (s *balanceService) Transfer(ctx context.Context, chatID, senderID, targetID int64, targetUsername string, delta int64) (*domain.TransferResult, error) {
	if senderID == targetID {
		return nil, domain.ErrSelfTransfer
	}

	absDelta := abs(delta)

	remaining, err := s.remaining(ctx, chatID, senderID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	if absDelta > remaining {
		return &domain.TransferResult{Remaining: remaining}, domain.ErrInsufficientLimit
	}

	target, err := s.users.GetByTelegramID(ctx, chatID, targetID)
	if errors.Is(err, domain.ErrUserNotFound) && targetUsername != "" {
		target, err = s.users.GetByUsername(ctx, chatID, targetUsername)
	}
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	opID := newOpID()
	now := nowUTC()
	postings := []domain.Posting{
		{ChatID: chatID, AccountID: target.TelegramID, Book: domain.BookScore, Amount: delta, OpID: opID, OpType: domain.OpTransfer, Counterparty: senderID, CreatedAt: now},
		{ChatID: chatID, AccountID: senderID, Book: domain.BookAllowance, Amount: absDelta, OpID: opID, OpType: domain.OpTransfer, Counterparty: target.TelegramID, CreatedAt: now},
	}
	if err := s.ledger.InsertPostings(ctx, postings); err != nil {
		slog.Error("transfer failed", "chat_id", chatID, "from", senderID, "to", target.TelegramID, "error", err)
		return nil, err
	}

	slog.Info("transfer", "chat_id", chatID, "from", senderID, "to", target.TelegramID, "delta", delta)
	return &domain.TransferResult{Target: target, Delta: delta}, nil
}

func (s *balanceService) TransferToAll(ctx context.Context, chatID, senderID int64, delta int64) (*domain.TransferAllResult, error) {
	recipients, err := s.users.ListUserIDsExcept(ctx, chatID, senderID)
	if err != nil {
		return nil, err
	}
	if len(recipients) == 0 {
		return nil, domain.ErrNoUsersInChat
	}

	absDelta := abs(delta)
	affected := len(recipients)
	totalCost := absDelta * int64(affected)

	remaining, err := s.remaining(ctx, chatID, senderID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	if totalCost > remaining {
		// Populate the result so the handler can show the shortfall.
		return &domain.TransferAllResult{Delta: delta, AffectedCnt: affected, TotalCost: totalCost, Remaining: remaining}, domain.ErrInsufficientLimit
	}

	opID := newOpID()
	now := nowUTC()
	postings := make([]domain.Posting, 0, affected+1)
	for _, rid := range recipients {
		postings = append(postings, domain.Posting{
			ChatID: chatID, AccountID: rid, Book: domain.BookScore, Amount: delta,
			OpID: opID, OpType: domain.OpTransferAll, Counterparty: senderID, CreatedAt: now,
		})
	}
	postings = append(postings, domain.Posting{
		ChatID: chatID, AccountID: senderID, Book: domain.BookAllowance, Amount: totalCost,
		OpID: opID, OpType: domain.OpTransferAll, Counterparty: 0, CreatedAt: now,
	})
	if err := s.ledger.InsertPostings(ctx, postings); err != nil {
		slog.Error("transfer to all failed", "chat_id", chatID, "from", senderID, "error", err)
		return nil, err
	}

	slog.Info("transfer to all", "chat_id", chatID, "from", senderID, "delta", delta, "affected", affected)
	return &domain.TransferAllResult{Delta: delta, AffectedCnt: affected, TotalCost: totalCost}, nil
}

func (s *balanceService) GetDailyBalances(ctx context.Context, chatID int64) ([]domain.DailyBalance, error) {
	users, err := s.users.ListChatUsers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	netSpent, err := s.ledger.ChatAllowanceSpentSince(ctx, chatID, startOfUTCDay())
	if err != nil {
		return nil, err
	}

	result := make([]domain.DailyBalance, 0, len(users))
	for _, u := range users {
		remaining := s.dailyLimit - netSpent[u.TelegramID]
		result = append(result, domain.DailyBalance{
			User:       u,
			Remaining:  remaining,
			DailyLimit: s.dailyLimit,
		})
	}
	return result, nil
}

func (s *balanceService) AdjustDailyLimit(ctx context.Context, chatID int64, adminUsername, targetUsername string, delta int64) (*domain.AdjustResult, error) {
	if !s.isSuperAdmin(adminUsername) {
		return nil, domain.ErrNotAuthorized
	}

	target, err := s.users.GetByUsername(ctx, chatID, targetUsername)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	netSpent, err := s.ledger.AllowanceSpentSince(ctx, chatID, target.TelegramID, startOfUTCDay())
	if err != nil {
		return nil, err
	}
	// Allowance can stack above the base limit (grants), but never drops below 0.
	oldRemaining := s.dailyLimit - netSpent
	newRemaining := oldRemaining + delta
	if newRemaining < 0 {
		newRemaining = 0
	}

	posting := domain.Posting{
		ChatID: chatID, AccountID: target.TelegramID, Book: domain.BookAllowance, Amount: oldRemaining - newRemaining,
		OpID: newOpID(), OpType: domain.OpAdminAdjust, Counterparty: 0, CreatedAt: nowUTC(),
	}
	if err := s.ledger.InsertPostings(ctx, []domain.Posting{posting}); err != nil {
		slog.Error("daily limit adjust failed", "chat_id", chatID, "target", targetUsername, "error", err)
		return nil, err
	}

	slog.Info("daily limit adjusted", "chat_id", chatID, "admin", adminUsername, "target", targetUsername, "delta", delta)
	return &domain.AdjustResult{OldRemaining: oldRemaining, NewRemaining: newRemaining}, nil
}

func (s *balanceService) DailyLimit() int64 {
	return s.dailyLimit
}

func (s *balanceService) remaining(ctx context.Context, chatID, telegramID int64) (int64, error) {
	netSpent, err := s.ledger.AllowanceSpentSince(ctx, chatID, telegramID, startOfUTCDay())
	if err != nil {
		return 0, err
	}
	return s.dailyLimit - netSpent, nil
}

func (s *balanceService) isSuperAdmin(username string) bool {
	for _, admin := range s.superAdmin {
		if strings.EqualFold(username, admin) {
			return true
		}
	}
	return false
}
