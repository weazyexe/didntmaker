package integration

import "testing"

func TestBetStats(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")
	env.register(bob, "bob")

	// applyBet calls ApplyResult directly, bypassing the once-per-day gate.
	env.applyBet(alice, 6) // win
	env.applyBet(alice, 6) // win
	env.applyBet(alice, 1) // lose

	won, lost, err := env.postingRepository.BetStats(env.ctx, env.chat, alice)
	if err != nil {
		t.Fatalf("bet stats: %v", err)
	}
	assertEq(t, won, 2)  // two wins
	assertEq(t, lost, 1) // one loss

	chatStats, err := env.betService.ChatBetStats(env.ctx, env.chat)
	if err != nil {
		t.Fatalf("chat bet stats: %v", err)
	}
	// Only alice played; bob is absent.
	assertEq(t, len(chatStats), 1)
	assertEq(t, chatStats[0].TelegramID, alice)
	assertEq(t, chatStats[0].Won, 2)
	assertEq(t, chatStats[0].Lost, 1)
}
