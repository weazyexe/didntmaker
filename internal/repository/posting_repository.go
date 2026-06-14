package repository

import (
	"context"
	"database/sql"
	"time"

	"weazyexe.dev/didntmaker/internal/database/gen"
	"weazyexe.dev/didntmaker/internal/domain"
)

type PostingRepository interface {
	InsertPostings(ctx context.Context, postings []domain.Posting) error
	Score(ctx context.Context, chatID, accountID int64) (int64, error)
	AllowanceSpentSince(ctx context.Context, chatID, accountID int64, since time.Time) (int64, error)
	ChatAllowanceSpentSince(ctx context.Context, chatID int64, since time.Time) (map[int64]int64, error)
	HasBetSince(ctx context.Context, chatID, accountID int64, since time.Time) (bool, error)
	Leaderboard(ctx context.Context, chatID int64) ([]domain.LeaderboardEntry, error)
	BetStats(ctx context.Context, chatID, accountID int64) (won, lost int64, err error)
	ChatBetStats(ctx context.Context, chatID int64) ([]domain.BetStatEntry, error)
}

type postingRepository struct {
	db      *sql.DB
	queries *gen.Queries
}

func NewPostingRepository(db *sql.DB) *postingRepository {
	return &postingRepository{db: db, queries: gen.New(db)}
}

func (r *postingRepository) InsertPostings(ctx context.Context, postings []domain.Posting) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for _, p := range postings {
		if err := qtx.InsertPosting(ctx, gen.InsertPostingParams{
			ChatID:       p.ChatID,
			AccountID:    p.AccountID,
			Book:         string(p.Book),
			Amount:       p.Amount,
			OpID:         p.OpID,
			OpType:       string(p.OpType),
			Counterparty: p.Counterparty,
			Metadata:     p.Metadata,
			CreatedAt:    p.CreatedAt,
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postingRepository) Score(ctx context.Context, chatID, accountID int64) (int64, error) {
	return r.queries.GetScore(ctx, gen.GetScoreParams{ChatID: chatID, AccountID: accountID})
}

func (r *postingRepository) AllowanceSpentSince(ctx context.Context, chatID, accountID int64, since time.Time) (int64, error) {
	return r.queries.GetAllowanceSpentSince(ctx, gen.GetAllowanceSpentSinceParams{
		ChatID:    chatID,
		AccountID: accountID,
		CreatedAt: since,
	})
}

func (r *postingRepository) ChatAllowanceSpentSince(ctx context.Context, chatID int64, since time.Time) (map[int64]int64, error) {
	rows, err := r.queries.GetChatAllowanceSpentSince(ctx, gen.GetChatAllowanceSpentSinceParams{
		ChatID:    chatID,
		CreatedAt: since,
	})
	if err != nil {
		return nil, err
	}

	spent := make(map[int64]int64, len(rows))
	for _, row := range rows {
		spent[row.AccountID] = row.Spent
	}
	return spent, nil
}

func (r *postingRepository) HasBetSince(ctx context.Context, chatID, accountID int64, since time.Time) (bool, error) {
	has, err := r.queries.HasBetSince(ctx, gen.HasBetSinceParams{
		ChatID:    chatID,
		AccountID: accountID,
		CreatedAt: since,
	})
	if err != nil {
		return false, err
	}
	return has != 0, nil
}

func (r *postingRepository) BetStats(ctx context.Context, chatID, accountID int64) (won, lost int64, err error) {
	row, err := r.queries.GetUserBetStats(ctx, gen.GetUserBetStatsParams{ChatID: chatID, AccountID: accountID})
	if err != nil {
		return 0, 0, err
	}
	return row.Won, row.Lost, nil
}

func (r *postingRepository) ChatBetStats(ctx context.Context, chatID int64) ([]domain.BetStatEntry, error) {
	rows, err := r.queries.GetChatBetStats(ctx, chatID)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.BetStatEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, domain.BetStatEntry{
			TelegramID: row.TelegramID,
			Username:   row.Username,
			FirstName:  row.FirstName,
			Won:        row.Won,
			Lost:       row.Lost,
		})
	}
	return entries, nil
}

func (r *postingRepository) Leaderboard(ctx context.Context, chatID int64) ([]domain.LeaderboardEntry, error) {
	rows, err := r.queries.GetLeaderboard(ctx, chatID)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.LeaderboardEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, domain.LeaderboardEntry{
			TelegramID: row.TelegramID,
			Username:   row.Username,
			FirstName:  row.FirstName,
			Score:      row.Score,
		})
	}
	return entries, nil
}
