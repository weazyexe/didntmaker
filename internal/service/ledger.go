package service

import (
	"time"

	"github.com/google/uuid"
)

// startOfUTCDay is the "today" boundary for daily allowance and bet checks —
// replaces the old daily-reset bookkeeping with a filter on created_at.
func startOfUTCDay() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func newOpID() string {
	return uuid.NewString()
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
