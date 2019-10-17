package types

// Bonus -
type Bonus struct {
	Multiplier uint `json:"multiplier"`
}

// NewBonus -
func NewBonus(amt uint) Bonus {
	return Bonus{
		Multiplier: amt,
	}
}
