package longy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

// GenesisState is the state that must be provided at genesis
type GenesisState struct {
	Attendees []types.Attendee
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Attendees: nil,
	}
}

// ValidateGenesis runs sanity checks `state`
func ValidateGenesis(state GenesisState) error {
	var seenIds map[string]bool
	var foundSuperUser bool
	for _, a := range state.Attendees {
		if seenIds[a.ID] {
			return fmt.Errorf("duplicate id: %s", a.ID)
		}
		seenIds[a.ID] = true

		if a.IsSuperUser() {
			if foundSuperUser {
				return fmt.Errorf("duplicate super user")
			} else if a.Address().Empty() {
				return fmt.Errorf("empty super user public key")
			}

			foundSuperUser = true
		} else if !a.Address().Empty() {
			return fmt.Errorf("normal attendee public keys must be empty on genesis")
		}
	}

	return nil
}

// InitGenesis will run module initialization using the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state GenesisState) {
	for _, a := range state.Attendees {
		if a.IsSuperUser() {
			k.SetMasterPublicKey(ctx, a.Address())
		} else {
			k.SetAttendee(ctx, a)
		}
	}
}
