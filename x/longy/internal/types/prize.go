package types

//nolint:golint
const (
	ClaimBadgeAwardPoints uint = 5
	//ScanAttendeeAwardPoints are the points given when you scan a regular attendee
	ScanAttendeeAwardPoints uint = 1
	//ScanSponsorAwardPoints are the points given when you scan a sponsor attendee
	ScanSponsorAwardPoints uint = 2

	//ShareAttendeeAwardPoints are the points when you share info with a regular attendee
	ShareAttendeeAwardPoints uint = 3
	//ShareSponsorAwardPoints are the points when you share info with a sponsor attendee
	ShareSponsorAwardPoints uint = 6

	//Tier is the rep needed to be in that tier level
	Tier1Rep uint = 30  //20
	Tier2Rep uint = 50  //30
	Tier3Rep uint = 75  //40
	Tier4Rep uint = 100 //50
	Tier5Rep uint = 125 //75
	Tier6Rep uint = 150 //100
	Tier7Rep uint = 175 //125
	Tier8Rep uint = 250 //170
	Tier9Rep uint = 600 //220

	//Quantity of prizes per tier
	Tier1Quantity uint = 150
	Tier2Quantity uint = 200
	Tier3Quantity uint = 300
	Tier4Quantity uint = 50
	Tier5Quantity uint = 150
	Tier6Quantity uint = 100
	Tier7Quantity uint = 100
	Tier8Quantity uint = 150
	Tier9Quantity uint = 1
)

//nolint:golint
const (
	Tier0 uint = iota
	Tier1
	Tier2
	Tier3
	Tier4
	Tier5
	Tier6
	Tier7
	Tier8
	Tier9
)

//Prize is the genesis type for the prizes
type Prize struct {
	Tier             uint   `json:"tier"`
	RepNeeded        uint   `json:"repNeeded"`
	PrizeText        string `json:"prizeText"`
	PrizeDescription string `json:"prizeDescription"`
	Quantity         uint   `json:"quantity"`
}

//Win is the type for the prizes that an attendee has earned playing the game
type Win struct {
	Tier    uint   `json:"tier"`
	Name    string `json:"name"`
	Claimed bool   `json:"claimed"`
}
