package types

// Bonus -
type Bonus struct {
	Multiplier int
}

// NewBonus -
func NewBonus(amt int) Bonus {
	return Bonus{
		Multiplier: amt,
	}
}
