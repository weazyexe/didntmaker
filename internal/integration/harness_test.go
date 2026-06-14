package integration

import (
	"context"
	"path/filepath"
	"testing"

	"weazyexe.dev/didntmaker/internal/database"
	"weazyexe.dev/didntmaker/internal/domain"
	"weazyexe.dev/didntmaker/internal/repository"
	"weazyexe.dev/didntmaker/internal/service"
)

const dailyLimit int64 = 1000

const (
	alice int64 = 100
	bob   int64 = 200
)

// env wires the whole stack against a fresh DB and exposes scenario-level
// helpers so tests read as actions + assertions, not plumbing.
type env struct {
	t                 *testing.T
	ctx               context.Context
	chat              int64
	usersRepository   repository.UserRepository
	postingRepository repository.PostingRepository
	userService       service.UserService
	balanceService    service.BalanceService
	betService        service.BetService
}

// setup boots a fresh DB with "alice" as the only super admin.
func setup(t *testing.T) *env {
	t.Helper()
	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	users := repository.NewUserRepository(db)
	postingRepository := repository.NewPostingRepository(db)
	return &env{
		t:                 t,
		ctx:               context.Background(),
		chat:              1,
		usersRepository:   users,
		postingRepository: postingRepository,
		userService:       service.NewUserService(users, postingRepository, dailyLimit),
		balanceService:    service.NewBalanceService(users, postingRepository, dailyLimit, []string{"alice"}),
		betService:        service.NewBetService(postingRepository, dailyLimit),
	}
}

func (e *env) register(id int64, name string) {
	e.t.Helper()
	if _, err := e.usersRepository.GetOrCreate(e.ctx, e.chat, id, name, name); err != nil {
		e.t.Fatalf("register %s: %v", name, err)
	}
}

func (e *env) transfer(from, to, delta int64) error {
	_, err := e.balanceService.Transfer(e.ctx, e.chat, from, to, "", delta)
	return err
}

func (e *env) mustTransfer(from, to, delta int64) {
	e.t.Helper()
	if err := e.transfer(from, to, delta); err != nil {
		e.t.Fatalf("transfer: %v", err)
	}
}

func (e *env) applyBet(id int64, dice int) {
	e.t.Helper()
	if _, err := e.betService.ApplyResult(e.ctx, e.chat, id, dice); err != nil {
		e.t.Fatalf("bet: %v", err)
	}
}

func (e *env) adjust(admin, target string, delta int64) (*domain.AdjustResult, error) {
	return e.balanceService.AdjustDailyLimit(e.ctx, e.chat, admin, target, delta)
}

func (e *env) score(id int64) int64 {
	e.t.Helper()
	score, err := e.postingRepository.Score(e.ctx, e.chat, id)
	if err != nil {
		e.t.Fatalf("score: %v", err)
	}
	return score
}

func (e *env) remaining(id int64) int64 {
	e.t.Helper()
	balances, err := e.balanceService.GetDailyBalances(e.ctx, e.chat)
	if err != nil {
		e.t.Fatalf("balances: %v", err)
	}
	for _, b := range balances {
		if b.User.TelegramID == id {
			return b.Remaining
		}
	}
	e.t.Fatalf("user %d not found in balances", id)
	return 0
}

func assertEq[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
