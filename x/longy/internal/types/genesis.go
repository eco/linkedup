package types

import (
	"fmt"
)

// GenesisState is the state that must be provided at genesis
type GenesisState struct {
	Attendees []Attendee
}

func ValidateGenesis(state GenesisState) error {
	var seenIds map[string]bool
	for _, a := range state.Attendees {
		if seenIds[a.ID] {
			return fmt.Errorf("duplicate id: %s", a.ID)
		}
		seenIds[a.ID] = true

		if !a.PublicKey.Empty() {
			return fmt.Errorf("attendee public keys must be empty on genesis")
		}
	}

	return nil
}
