package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Attendee encapsulates attendee information
type Attendee struct {
	ID        string
	PublicKey sdk.AccAddress
}
