package bot

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
	"weazyexe.dev/didntmaker/internal/i18n"
)

func TestRateLimiterBurstThenDrop(t *testing.T) {
	var aliceID int64 = 67
	var bobID int64 = 52

	rl := newRateLimiter(rate.Every(time.Hour), 3, i18n.Get(i18n.Default()))

	for i := 0; i < 3; i++ {
		if !rl.allow(aliceID) {
			t.Fatalf("call %d within burst should be allowed", i)
		}
	}
	if rl.allow(aliceID) {
		t.Fatal("4th call over burst should be dropped")
	}
	// separate user has its own bucket
	if !rl.allow(bobID) {
		t.Fatal("different user should not be rate limited")
	}
}
