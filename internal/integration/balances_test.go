package integration

import "testing"

func TestBalancesBetAvailable(t *testing.T) {
	env := setup(t)
	env.register(alice, "alice")
	env.register(bob, "bob")

	env.applyBet(alice, 6) // alice крутит сегодня

	balances, err := env.balanceService.GetDailyBalances(env.ctx, env.chat)
	if err != nil {
		t.Fatalf("balances: %v", err)
	}

	avail := map[int64]bool{}
	for _, b := range balances {
		avail[b.User.TelegramID] = b.BetAvailable
	}
	assertEq(t, avail[alice], false) // уже крутил
	assertEq(t, avail[bob], true)    // ещё доступна
}
