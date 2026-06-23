package domain

type UserStats struct {
	User           User
	Score          int64
	DailyRemaining int64
	DailyLimit     int64
	Won            int64
	Lost           int64
	BetAvailable   bool

	WeekDelta  int64
	MonthDelta int64

	WorstDayMinus int64 // исторический максимум минусов, полученных за один день
	BestDayPlus   int64 // исторический максимум плюсов, полученных за один день

	Fan      *Counterparty // кто тебя больше всех плюсует
	Hater    *Counterparty // кто больше всех минусует
	Favorite *Counterparty // кого ты больше всех плюсуешь
	Victim   *Counterparty // кого больше всех минусуешь
}

// Counterparty is the other participant in a stats line; Amount is a positive magnitude.
type Counterparty struct {
	Username  string
	FirstName string
	Amount    int64
}

// CounterpartyAgg is one grouped row: a participant with their total plus/minus.
type CounterpartyAgg struct {
	Username  string
	FirstName string
	Plus      int64
	Minus     int64
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
