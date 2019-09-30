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
	rep     int
}

// NewAttendee is the constructor for `Attendee`. New attendees default to 0 rep
func NewAttendee(id string) Attendee {
	addr := util.IDToAddress(id)

	return Attendee{
		id:      id,
		address: addr,
		rep:     0,
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
