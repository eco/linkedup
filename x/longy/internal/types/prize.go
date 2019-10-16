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
	Tier1Rep uint = 20
	Tier2Rep uint = 50
	Tier3Rep uint = 100
	Tier4Rep uint = 350
	Tier5Rep uint = 400
)

//nolint:golint
const (
	Tier0 uint = iota
	Tier1
	Tier2
	Tier3
	Tier4
	Tier5
)

//Prize is the genesis type for the prizes
type Prize struct {
	Tier      uint   `json:"tier"`
	RepNeeded uint   `json:"repNeeded"`
	PrizeText string `json:"prizeText"`
	Quantity  int    `json:"quantity"`
}

//Win is the type for the prizes that an attendee has earned playing the game
type Win struct {
	Tier    uint   `json:"tier"`
	Name    string `json:"name"`
	Claimed bool   `json:"claimed"`
}
