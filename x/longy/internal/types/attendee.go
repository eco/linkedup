package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Attendee encapsulates attendee information
type Attendee struct {
	ID        string
	PublicKey sdk.AccAddress
}

// GetID returns the attendee's identifier
func (a Attendee) GetID() string {
	return a.ID
}

// PublicKey returns the public key associated with the attendee
func (a Attendee) Address() sdk.AccAddress {
	return a.PublicKey
}
