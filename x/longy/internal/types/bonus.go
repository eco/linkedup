package types

import (
	"strconv"
)

// Bonus -
type Bonus struct {
	Multiplier string `json:"multiplier"`
}

// NewBonus -
func NewBonus(amt string) Bonus {
	return Bonus{
		Multiplier: amt,
	}
}

// GetMultiplier - will return 1 if it fails to convert
func (b *Bonus) GetMultiplier() float64 {
	num, err := strconv.ParseFloat(b.Multiplier, 64)
	if err != nil {
		// This should NEVER happen. But if it does....
		return 1
	}

	return num
}
