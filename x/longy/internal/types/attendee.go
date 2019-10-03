package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
)

// GenesisAttendee is all the information needed to create an
// attendee at genesis
type GenesisAttendee struct {
	ID string `json:"id"`
}

// Attendee encapsulates attendee information
type Attendee struct {
	ID      string
	Address sdk.AccAddress

	Commitment util.Commitment
	Claimed    bool

	Rep int
}

// NewAttendee is the constructor for `Attendee`. New attendees default to 0 rep
// and is unclaimed
func NewAttendee(id string) Attendee {
	addr := util.IDToAddress(id)

	return Attendee{
		ID:      id,
		Address: addr,

		Commitment: nil,
		Claimed:    false,

		Rep: 0,
	}
}

// NewAttendeeFromGenesis will create an `Attendee` from `GenesisAttendee`
func NewAttendeeFromGenesis(ga GenesisAttendee) Attendee {
	return NewAttendee(ga.ID)
}

// GetID returns the attendee identifier
func (a Attendee) GetID() string {
	return a.ID
}

// GetAddress returns the deterministic address associated with the attendee
func (a Attendee) GetAddress() sdk.AccAddress {
	return a.Address
}

// GetRep returns the attendee's current rep
func (a Attendee) GetRep() int {
	return a.Rep
}

// CurrentCommitment returns the current commitment associated with this attendee
func (a Attendee) CurrentCommitment() util.Commitment {
	return a.Commitment
}

// SetCommitment will set the claim hash for this attendee
func (a *Attendee) SetCommitment(commitment util.Commitment) {
	bz := make([]byte, commitment.Len())
	copy(bz[:], commitment.Bytes())
	a.Commitment = bz
}

// ResetCommitment will reset this attendee's commitment to nil
func (a *Attendee) ResetCommitment() {
	a.Commitment = nil
}

// IsClaimed indicates if this attendee is claimed
func (a Attendee) IsClaimed() bool {
	return a.Claimed
}

// SetClaimed will mark this attendee as claimed
func (a *Attendee) SetClaimed() {
	a.Claimed = true
}

// SetUnclaimed will mark this attendee as unclaimed
func (a *Attendee) SetUnclaimed() {
	a.Claimed = false
}
