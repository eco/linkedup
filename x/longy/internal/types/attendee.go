package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Attendee is the structure of the attendees
type Attendee struct {
	ID      string         `json:"id"`
	Address sdk.AccAddress `json:"address"`
}

// GetIDBytes returns the ID for an attendee, also known as the badgeID from eventbright
func (e *Attendee) GetIDBytes() []byte {
	return []byte(e.ID)
}
