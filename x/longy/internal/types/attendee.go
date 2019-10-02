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
	id      string
	address sdk.AccAddress

	commitment util.Commitment
	claimed    bool

	rep int
}

// NewAttendee is the constructor for `Attendee`. New attendees default to 0 rep
// and is unclaimed
func NewAttendee(id string) Attendee {
	addr := util.IDToAddress(id)

	return Attendee{
		id:      id,
		address: addr,

		commitment: util.Commitment{},
		claimed:    false,

		rep: 0,
	}
}

// NewAttendeeFromGenesis will create an `Attendee` from `GenesisAttendee`
func NewAttendeeFromGenesis(ga GenesisAttendee) Attendee {
	return NewAttendee(ga.ID)
}

// ID returns the attendee identifier
func (a Attendee) ID() string {
	return a.id
}

// Address returns the deterministic address associated with the attendee
func (a Attendee) Address() sdk.AccAddress {
	return a.address
}

// Rep returns the attendee's current rep
func (a Attendee) Rep() int {
	return a.rep
}

// CurrentCommitment returns the current commitment associated with this attendee
func (a Attendee) CurrentCommitment() util.Commitment {
	return a.commitment
}

// SetCommitment will set the claim hash for this attendee
func (a *Attendee) SetCommitment(commitment util.Commitment) {
	copy(a.commitment[:], commitment)
}

// ResetCommitment will reset this attendee's commitment to nil
func (a *Attendee) ResetCommitment() {
	a.commitment = util.Commitment{}
}

// IsClaimed indicates if this attendee is claimed
func (a Attendee) IsClaimed() bool {
	return a.claimed
}

// SetClaim will mark this attendee as claimed
func (a *Attendee) SetClaimed() {
	a.claimed = true
}

// SetUnclaim will mark this attendee as unclaimed
func (a *Attendee) SetUnclaimed() {
	a.claimed = false
}
