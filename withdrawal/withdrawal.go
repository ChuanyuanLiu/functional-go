package withdrawal

import (
	"sort"
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

// monthly withdrawal upper-bound to fee, e.g. 0-100 has no fees
var tiers = map[float32]int{
	100:  0,
	1000: 10,
	1500: 50,
	2000: 100,
}

const maxWithdrawalFee = 200

func CalculateFee(now time.Time, pastWithdrawals []Withdrawal, amount float32) int {
	thisMonth := time.Date(now.UTC().Year(), now.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)

	completedWithdrawal := func(w Withdrawal, _ int) bool {
		return w.Date.After(thisMonth) && w.Status != Failure
	}

	getAmount := func(w Withdrawal, _ int) float32 {
		return w.Amount
	}

	total := lo.Sum(lo.Map(lo.Filter(pastWithdrawals, completedWithdrawal), getAmount))

	fee := 0

	sorted_tiers := lo.Entries(tiers)
	sort.Slice(sorted_tiers, func(i, j int) bool {
		return sorted_tiers[i].Key < sorted_tiers[j].Key
	})

	for amount > 0 {
		sorted_tiers = lo.DropWhile(sorted_tiers, func(x lo.Entry[float32, int]) bool {
			return total >= x.Key
		})
		if len(sorted_tiers) == 0 {
			fee += maxWithdrawalFee
			break
		}
		amount -= sorted_tiers[0].Key - total
		total = sorted_tiers[0].Key
		fee += sorted_tiers[0].Value
	}

	return fee
}
