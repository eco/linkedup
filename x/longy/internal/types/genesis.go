package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisAttendee is the attendee structure in the genesis file
type GenesisAttendee struct {
	ID string `json:"id"`
	//Profile GenesisProfile `json:"profile"`   //gets the full info of the account
}

// GenesisAttendees is the full array of attendees to initialize
type GenesisAttendees []GenesisAttendee

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
