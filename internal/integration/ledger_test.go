package integration

import (
	"testing"

	"weazyexe.dev/didntmaker/internal/domain"
)

func TestTransfer(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")
	env.register(bob, "bob")

	env.mustTransfer(alice, bob, -100)

	assertEq(t, env.score(bob), -100)      // receiver's score moves
	assertEq(t, env.score(alice), 0)       // sender's score does not
	assertEq(t, env.remaining(alice), 900) // sender spends allowance
	assertEq(t, env.remaining(bob), 1000)  // sender spends allowance
}

func TestTransferLimitedByDailyAllowance(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")
	env.register(bob, "bob")

	env.mustTransfer(alice, bob, 600) // remaining 400
	if err := env.transfer(alice, bob, 600); err != domain.ErrInsufficientLimit {
		t.Fatalf("err = %v, want ErrInsufficientLimit", err)
	}
}

func TestTransferInsufficientReportsRemaining(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")
	env.register(bob, "bob")

	env.mustTransfer(alice, bob, -700) // remaining 300

	res, err := env.balanceService.Transfer(env.ctx, env.chat, alice, bob, "", -500)
	if err != domain.ErrInsufficientLimit {
		t.Fatalf("err = %v, want ErrInsufficientLimit", err)
	}
	assertEq(t, res.Remaining, 300) // real remaining, not a hardcoded 0
}

func TestSingleTransferCanExceedThousand(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")
	env.register(bob, "bob")

	env.applyBet(alice, 6)              // win => allowance stacks to 2000
	env.mustTransfer(alice, bob, -1500) // a single transfer above the old per-tx cap

	assertEq(t, env.score(bob), -1500)
	assertEq(t, env.remaining(alice), 500)
}

func TestBetLossLocksForTheDay(t *testing.T) {
	env := setup(t)
	env.register(bob, "bob")

	if err := env.betService.CanBet(env.ctx, env.chat, bob); err != nil {
		t.Fatalf("canbet before: %v", err)
	}

	env.applyBet(bob, 1) // dice 1-3 => lose
	assertEq(t, env.score(bob), -dailyLimit)

	if err := env.betService.CanBet(env.ctx, env.chat, bob); err != domain.ErrBetAlreadyUsed {
		t.Fatalf("canbet after = %v, want ErrBetAlreadyUsed", err)
	}
}

func TestBetWinStacksAllowance(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")

	assertEq(t, env.remaining(alice), dailyLimit) // 1000, nothing spent

	env.applyBet(alice, 6) // dice 4-6 => win, +1000 allowance

	assertEq(t, env.remaining(alice), 2*dailyLimit) // stacks above the base limit
}

func TestAdminAdjust(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice") // super admin
	env.register(bob, "bob")

	env.mustTransfer(bob, alice, 100) // bob spends 100 -> remaining 900

	res, err := env.adjust("alice", "bob", 200)
	if err != nil {
		t.Fatalf("adjust: %v", err)
	}
	assertEq(t, res.OldRemaining, 900)
	assertEq(t, res.NewRemaining, 1100) // stacks above the base limit

	if _, err := env.adjust("bob", "alice", 100); err != domain.ErrNotAuthorized {
		t.Fatalf("non-admin adjust = %v, want ErrNotAuthorized", err)
	}
}

func TestAdminAdjustFloorsAllowanceAtZero(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice") // super admin
	env.register(bob, "bob")

	res, err := env.adjust("alice", "bob", -10_000_000)
	if err != nil {
		t.Fatalf("adjust: %v", err)
	}
	assertEq(t, res.NewRemaining, 0) // allowance never goes below 0
	assertEq(t, env.remaining(bob), 0)
}
