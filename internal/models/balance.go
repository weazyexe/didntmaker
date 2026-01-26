package models

type TransferResult struct {
	Target     *User
	Delta      int64
	OldBalance int64
}

type TransferAllResult struct {
	Delta       int64
	AffectedCnt int
	TotalCost   int64
}

type DailyBalance struct {
	User       *User
	Remaining  int64
	DailyLimit int64
}

type AdjustResult struct {
	OldRemaining int64
	NewRemaining int64
}
