package integration

import "testing"

// Confirms the day-grouping SQL (date(created_at)) survives the real driver and
// that worst/best day totals aggregate same-day minuses/pluses correctly.
func TestExtremeDaysReceived(t *testing.T) {
	e := setup(t)
	e.register(alice, "alice")
	e.register(bob, "bob")

	const carol int64 = 300
	e.register(carol, "carol")

	// alice gets minuses from two people and a plus, all the same day
	e.mustTransfer(bob, alice, -300)
	e.mustTransfer(carol, alice, -200)
	e.mustTransfer(bob, alice, +100)

	stats, err := e.userService.GetStats(e.ctx, e.chat, alice, "", "alice")
	if err != nil {
		t.Fatalf("stats: %v", err)
	}

	assertEq(t, stats.WorstDayMinus, int64(500)) // -300 + -200, as positive magnitude
	assertEq(t, stats.BestDayPlus, int64(100))
}
