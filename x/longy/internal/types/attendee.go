package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
)

// GenesisAttendees is the full array of attendees to initialize
type GenesisAttendees []GenesisAttendee

// GenesisAttendee is the attendee structure in the genesis file
type GenesisAttendee struct {
	ID string `json:"id"`
	//Profile GenesisProfile `json:"profile"`   //gets the full info of the account
}

// GenesisProfile is the profile of the attendee from eventbrite
type GenesisProfile struct {
	Name     string `json:"name"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}

// GenesisService is the genesis type for the re-key service and its account address
type GenesisService struct {
	Address sdk.AccAddress `json:"address"`
}

// Attendee encapsulates attendee information
type Attendee struct {
	ID      string
	Address sdk.AccAddress

	Commitment       util.Commitment
	Claimed          bool
	FirstTimeClaimed bool

	Rep uint
}

// NewAttendee is the constructor for `Attendee`. New attendees default to 0 rep
// and is unclaimed
func NewAttendee(id string) Attendee {
	addr := util.IDToAddress(id)

	return Attendee{
		ID:      id,
		Address: addr,

		Commitment:       nil,
		Claimed:          false,
		FirstTimeClaimed: false,

		Rep: 0,
	}
}

// NewAttendeeFromGenesis will create an `Attendee` from `GenesisAttendee`
func NewAttendeeFromGenesis(ga GenesisAttendee) Attendee {
	return NewAttendee(ga.ID)
}

// GetID returns the attendee identifier
//nolint:gocritic
func (a Attendee) GetID() string {
	return a.ID
}

// GetAddress returns the deterministic address associated with the attendee
//nolint:gocritic
func (a Attendee) GetAddress() sdk.AccAddress {
	return a.Address
}

// GetRep returns the attendee's current rep
//nolint:gocritic
func (a Attendee) GetRep() uint {
	return a.Rep
}

// AddRep will add rep to the attendee {
func (a *Attendee) AddRep(r uint) {
	a.Rep += r
}

// CurrentCommitment returns the current commitment associated with this attendee
//nolint:gocritic
func (a Attendee) CurrentCommitment() util.Commitment {
	return a.Commitment
}

// SetCommitment will set the claim hash for this attendee
//nolint:gocritic
func (a *Attendee) SetCommitment(commitment util.Commitment) {
	bz := make([]byte, commitment.Len())
	copy(bz, commitment.Bytes())
	a.Commitment = bz
}

// ResetCommitment will reset this attendee's commitment to nil
func (a *Attendee) ResetCommitment() {
	a.Commitment = nil
}

// IsClaimed indicates if this attendee is claimed
//nolint:gocritic
func (a Attendee) IsClaimed() bool {
	return a.Claimed
}

// SetClaimed will mark this attendee as claimed
func (a *Attendee) SetClaimed() {
	a.Claimed = true
	if !a.FirstTimeClaimed {
		a.FirstTimeClaimed = true
	}
}
