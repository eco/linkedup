package longy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

// GenesisState is the genesis struct for the longy module
type GenesisState struct {
	Service   GenesisService   `json:"service"`
	Attendees GenesisAttendees `json:"attendees"`
}

// DefaultGenesisState returns the default genesis struct for the longy module
func DefaultGenesisState() GenesisState {
	return GenesisState{Service: GenesisService{}, Attendees: GenesisAttendees{}}
}

// ValidateGenesis validates that the passed genesis state is valid
func ValidateGenesis(data GenesisState) error {
	if data.Service.Address.Empty() {
		return types.ErrGenesisServiceAddressEmpty("Re-Key Service address must be set")
	}

	if data.Attendees == nil {
		return types.ErrGenesisAttendeesEmpty("Genesis attendees cannot be empty")
	}

	var seenIds = make(map[string]bool)
	for _, a := range data.Attendees {
		if seenIds[a.ID] {
			return fmt.Errorf("duplicate id: %s", a.ID)
		}
		seenIds[a.ID] = true
	}
	return nil
}

// InitGenesis will run module initialization using the genesis state
//nolint:gocritic
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state GenesisState) {
	// set the master public key
	k.SetMasterPublicKey(ctx, state.Service.Address)

	// create and set of all the attendees and cosmos accounts
	accountKeeper := k.AccountKeeper()
	for _, a := range state.Attendees {
		attendee := types.NewAttendeeFromGenesis(a)
		k.SetAttendee(ctx, attendee)

		account := accountKeeper.NewAccountWithAddress(ctx, attendee.GetAddress())
		accountKeeper.SetAccount(ctx, account)
	}
}
