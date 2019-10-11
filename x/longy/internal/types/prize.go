package types

//nolint:golint
const (
	//ScanAttendeeAwardPoints are the points given when you scan a regular attendee
	ScanAttendeeAwardPoints uint = 1
	//ScanSponsorAwardPoints are the points given when you scan a sponsor attendee
	ScanSponsorAwardPoints uint = 2

	//ShareAttendeeAwardPoints are the points when you share info with a regular attendee
	ShareAttendeeAwardPoints uint = 3
	//ShareSponsorAwardPoints are the points when you share info with a sponsor attendee
	ShareSponsorAwardPoints uint = 6

	//Tier is the rep needed to be in that tier level
	Tier1Rep = 100
	Tier2Rep = 200
	Tier3Rep = 300
	Tier4Rep = 350
	Tier5Rep = 400
)

//nolint:golint
const (
	Tier0 = iota
	Tier1
	Tier2
	Tier3
	Tier4
	Tier5
)

//Prize is the genesis type for the prizes
type Prize struct {
	Tier      int    `json:"tier"`
	RepNeeded int    `json:"repNeeded"`
	PrizeText string `json:"prizeText"`
	Quantity  int    `json:"quantity"`
}
