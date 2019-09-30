package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Attendee encapsulates attendee information
type Attendee struct {
	ID        string
	PublicKey sdk.AccAddress

	isSuperUser bool
	rep         int
}

// NewAttendee is the constructor for `Attendee`. New attendee's default to 0 rep
func NewAttendee(id string, publicKey sdk.AccAddress, isSuperUser bool) Attendee {
	return Attendee{
		ID:          id,
		PublicKey:   publicKey,
		isSuperUser: isSuperUser,

		rep: 0,
	}
}

// GetID returns the attendee's identifier
func (a Attendee) GetID() string {
	return a.ID
}

// Address returns the public key associated with the attendee
func (a Attendee) Address() sdk.AccAddress {
	return a.PublicKey
}

// IsSuperUser is the indicator for the master account
func (a Attendee) IsSuperUser() bool {
	return a.isSuperUser
}

// Rep returns the attendee's current rep
func (a Attendee) Rep() int {
	return a.rep
}

// AddRep will increment `a`'s rep by `change` amount. An unsigned integer must
// be passed in to enforce a positive change
func (a *Attendee) AddRep(change uint) {
	a.rep += int(change)
}
