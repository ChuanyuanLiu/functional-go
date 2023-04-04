package withdrawal

import (
	"time"

	"github.com/samber/lo"
)

type status string

const (
	Success    status = "success"
	Failure    status = "failed"
	InProgress status = "inprogress"
)

type Withdrawal struct {
	UserID int
	Amount float32
	Status status
	Date   time.Time
}

// monthly withdrawal lower-bound to fee, e.g. 0-99 has no fees
var tiers = map[float32]int{
	0:    0,
	100:  10,
	1000: 50,
	1500: 100,
}

var highestTierAmount float32 = 5000
var highestTierPercent float32 = 0.05

func CalculateFee(now time.Time, pastWithdrawals []Withdrawal, amount float32) int {
	thisMonth := time.Date(now.UTC().Year(), now.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)

	completedWithdrawal := func(w Withdrawal, _ int) bool {
		return w.Date.After(thisMonth) && w.Status != Failure
	}

	getAmount := func(w Withdrawal, _ int) float32 {
		return w.Amount
	}

	alreadyWithdrawn := lo.Sum(lo.Map(lo.Filter(pastWithdrawals, completedWithdrawal), getAmount))

	tiers := lo.Filter(lo.Entries(tiers), func(x lo.Entry[float32, int], _ int) bool {
		return alreadyWithdrawn < x.Key && x.Key <= alreadyWithdrawn+amount
	})

	fee := lo.Sum(lo.Map(tiers, func(x lo.Entry[float32, int], _ int) int {
		return x.Value
	}))

	// percentage rate
	if alreadyWithdrawn+amount >= highestTierAmount {
		fee += int((amount + alreadyWithdrawn - highestTierAmount) * highestTierPercent)
	}

	return fee
}
