package types

import "time"

const (
	//LeaderBoardCount is the total number of attendees on the board
	LeaderBoardCount = 30
	//LeaderBoardTier1Prize is the tier 1 prize pool in USD
	LeaderBoardTier1Prize = 5000
	//LeaderBoardTier2Prize is the tier 2 prize pool in USD
	LeaderBoardTier2Prize = 2000
	//LeaderBoardTier1Count is the number of people in the first tier
	LeaderBoardTier1Count = 10
)

//Tier is a prize tier in the leader  board
type Tier struct {
	PrizeAmount int        `json:"prizeAmount"`
	Attendees   []Attendee `json:"attendees"`
}

//LeaderBoard is the leader board struct
type LeaderBoard struct {
	TotalCount int  `json:"totalCount"`
	Tier1      Tier `json:"tier1"`
	Tier2      Tier `json:"tier2"`
	Time       time.Time
}

//NewLeaderBoard returns an initialized leader board with some constants
func NewLeaderBoard(count int, top []Attendee) *LeaderBoard { // test this
	var first, second []Attendee
	topLen := len(top)
	if topLen > LeaderBoardCount {
		topLen = LeaderBoardCount
	}

	if topLen >= LeaderBoardTier1Count {
		first = make([]Attendee, LeaderBoardTier1Count)
		copy(first, top[:LeaderBoardTier1Count])

		secondLen := topLen - LeaderBoardTier1Count
		second = make([]Attendee, secondLen)
		copy(second, top[LeaderBoardTier1Count:])
	} else {
		first = top
	}

	return &LeaderBoard{
		TotalCount: count,
		Tier1: Tier{
			PrizeAmount: LeaderBoardTier1Prize,
			Attendees:   first,
		},
		Tier2: Tier{
			PrizeAmount: LeaderBoardTier2Prize,
			Attendees:   second,
		},
	}
}
