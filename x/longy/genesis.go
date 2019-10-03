package longy

import "github.com/eco/longy/x/longy/errors"

//GenesisState is the genesis struct for the longy module
type GenesisState struct {
	Service   GenesisService   `json:"service"`
	Attendees GenesisAttendees `json:"attendees"`
}

//DefaultGenesisState returns the default genesis struct for the longy module
func DefaultGenesisState() GenesisState {
	return GenesisState{Service: GenesisService{}, Attendees: GenesisAttendees{}}
}

//ValidateGenesis validates that the passed genesis state is valid
func ValidateGenesis(data GenesisState) error {
	if data.Service.Address.Empty() {
		return errors.ErrGenesisServiceAddressEmpty("Re-Key Service address must be set")
	}

	if data.Attendees == nil {
		return errors.ErrGenesisAttendeesEmpty("Genesis attendees cannot be empty")
	}
	return nil
}
