package types

// Bonus -
type Bonus struct {
	Multiplier uint
}

// NewBonus -
func NewBonus(amt uint) Bonus {
	return Bonus{
		Multiplier: amt,
	}
}
