package genesis

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/spf13/cobra"
)

// AddGenesisPrizesCmd returns add-genesis-prizes cobra Command. Allows users to add the list of prizes
// to the chain and the number available
func AddGenesisPrizesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-prizes",
		Short: "Add genesis prizes to genesis.json",
		RunE: func(_ *cobra.Command, args []string) error {
			return addGenesisPrizes(ctx, cdc)
		},
	}

	return cmd
}

// addGenesisPrizes adds the conference prizes to the genesis file
func addGenesisPrizes(ctx *server.Context, cdc *codec.Codec) error {
	appState, genDoc, genFile, err := getGenesisState(ctx, cdc)
	if err != nil {
		return err
	}
	genesisState := buildPrizeGenesisState(appState, cdc)

	return updateGenesisState(cdc, genesisState, appState, genDoc, genFile)
}

func buildPrizeGenesisState(appState map[string]json.RawMessage, cdc *codec.Codec) longy.GenesisState {
	var (
		genesisState longy.GenesisState
	)

	// un-marshal the current state of the genesis object
	cdc.MustUnmarshalJSON(appState[longy.ModuleName], &genesisState)

	//get the prizes
	genesisState.Prizes = types.GetGenesisPrizes()

	return genesisState
}
