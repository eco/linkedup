package longy

//GenesisState is the genesis struct for the longy module
type GenesisState struct {
}

//DefaultGenesisState returns the default genesis struct for the longy module
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

//ValidateGenesis validates that the passed genesis state is valid
func ValidateGenesis(data GenesisState) error {
	return nil
}
