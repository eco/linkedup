package longy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

// GenesisState is the state that must be provided at genesis
type GenesisState struct {
	MasterKey sdk.AccAddress
	Attendees []types.GenesisAttendee
}

// DefaultGenesisState is an empty `GenesisState`
func DefaultGenesisState() GenesisState {
	return GenesisState{
		MasterKey: nil,
		Attendees: nil,
	}
}

// ValidateGenesis runs sanity checks `state`
func ValidateGenesis(state GenesisState) error {
	if state.MasterKey.Empty() {
		return fmt.Errorf("empty master key")
	}

	var seenIds map[string]bool
	for _, a := range state.Attendees {
		if seenIds[a.ID] {
			return fmt.Errorf("duplicate id: %s", a.ID)
		}
		seenIds[a.ID] = true
	}

	return nil
}

// InitGenesis will run module initialization using the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state GenesisState) {
	// set the master public key
	k.SetMasterPublicKey(ctx, state.MasterKey)

	// create and set all the attendees
	for _, a := range state.Attendees {
		attendee := types.NewAttendeeFromGenesis(a)
		k.SetAttendee(ctx, attendee)
	}
}
