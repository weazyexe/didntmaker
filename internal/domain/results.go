package domain

type UserStats struct {
	User           User
	Score          int64
	DailyRemaining int64
	DailyLimit     int64
	Won            int64
	Lost           int64
	BetAvailable   bool
}

type LeaderboardEntry struct {
	TelegramID int64
	Username   string
	FirstName  string
	Score      int64
}

type BetStatEntry struct {
	TelegramID int64
	Username   string
	FirstName  string
	Won        int64
	Lost       int64
}

type TransferResult struct {
	Target    User
	Delta     int64
	Remaining int64
}

type TransferAllResult struct {
	Delta       int64
	AffectedCnt int
	TotalCost   int64
	Remaining   int64
}

type DailyBalance struct {
	User         User
	Remaining    int64
	DailyLimit   int64
	BetAvailable bool
}

type AdjustResult struct {
	OldRemaining int64
	NewRemaining int64
}

type BetResult struct {
	DiceValue  int
	Won        bool
	DailyLimit int64
}
