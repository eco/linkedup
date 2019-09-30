package types

// Reward encapsulates the information for rep rewards
type Reward struct {
	kind string
	rep  uint
}

func NewReward(kind string, rep uint) Reward {
	return Reward{
		kind: kind,
		rep:  rep,
	}
}

// Kind returns the reward type
func (r Reward) Kind() string {
	return r.kind
}

// Rep returns the amount of rep this reward is worth
func (r Reward) Rep() uint {
	return r.rep
}
