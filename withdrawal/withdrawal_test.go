package withdrawal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFee(t *testing.T) {

	var pastWithdrawals = []Withdrawal{
		{
			UserID: 1,
			Amount: 30,
			Status: Success,
			Date:   time.Now(),
		},
		{
			UserID: 1,
			Amount: 40,
			Status: Failure,
			Date:   time.Now(),
		},
		{
			UserID: 1,
			Amount: 170,
			Status: InProgress,
			Date:   time.Now(),
		},
	}

	t.Run("Withdrawal amount passes 2 tiers ", func(t *testing.T) {
		fee := CalculateFee(time.Now(), pastWithdrawals, 1700)
		assert.Equal(t, 150, fee)
	})

	t.Run("Reaches the maximum tier", func(t *testing.T) {
		fee := CalculateFee(time.Now(), pastWithdrawals, 9000)
		assert.Equal(t, 360, fee)
	})

	var pastWithdrawals2 = []Withdrawal{
		{
			UserID: 1,
			Amount: 5000,
			Status: Success,
			Date:   time.Now(),
		},
	}

	t.Run("Starts at the maximum tier", func(t *testing.T) {
		fee := CalculateFee(time.Now(), pastWithdrawals2, 310)
		assert.Equal(t, 15, fee)
	})
}
