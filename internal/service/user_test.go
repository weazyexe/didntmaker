package service

import (
	"testing"

	"weazyexe.dev/didntmaker/internal/domain"
)

func TestTopPlusAndMinus(t *testing.T) {
	aggs := []domain.CounterpartyAgg{
		{Username: "masha", Plus: 2300, Minus: 100},
		{Username: "petya", Plus: 50, Minus: 1800},
		{Username: "vasya", Plus: 1200, Minus: 0},
	}

	if fan := topPlus(aggs); fan == nil || fan.Username != "masha" || fan.Amount != 2300 {
		t.Fatalf("topPlus = %+v, want masha/2300", fan)
	}
	if hater := topMinus(aggs); hater == nil || hater.Username != "petya" || hater.Amount != 1800 {
		t.Fatalf("topMinus = %+v, want petya/1800", hater)
	}
}

func TestTopEmptyAndZero(t *testing.T) {
	if topPlus(nil) != nil || topMinus(nil) != nil {
		t.Fatal("empty input should give nil")
	}
	// only-minus participants -> no fan; only-plus -> no hater
	onlyMinus := []domain.CounterpartyAgg{{Username: "x", Plus: 0, Minus: 500}}
	if topPlus(onlyMinus) != nil {
		t.Fatal("no positive totals should give nil fan")
	}
	if topMinus(onlyMinus) == nil {
		t.Fatal("a minus total should give a hater")
	}
}
