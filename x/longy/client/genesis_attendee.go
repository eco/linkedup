package client

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// AddGenesisAttendeesCmd returns add-genesis-attendees cobra Command. Allows users to add the list of attendees
// to the chain by their eventbrite id
func AddGenesisAttendeesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-attendees",
		Short: "Add genesis attendees to genesis.json",
		RunE: func(_ *cobra.Command, args []string) error {
			return addGenesisAttendees(ctx, cdc)
		},
	}

	return cmd
}

// AddGenesisAttendees adds the attendees and the service account to the genesis file under the longy key
func addGenesisAttendees(ctx *server.Context, cdc *codec.Codec) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))

	// retrieve the app state
	genFile := config.GenesisFile()
	appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
	if err != nil {
		return err
	}

	genesisState, err := buildGenesisState(appState, cdc)
	if err != nil {
		return err
	}

	genesisStateBz := cdc.MustMarshalJSON(genesisState)
	appState[longy.ModuleName] = genesisStateBz

	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	// export app state
	genDoc.AppState = appStateJSON

	return genutil.ExportGenesisFile(genDoc, genFile)
}

// BuildGenesisState builds the genesis state for the longy module
func buildGenesisState(appState map[string]json.RawMessage, cdc *codec.Codec) (longy.GenesisState, sdk.Error) {
	var (
		genesisState longy.GenesisState
		err          sdk.Error
	)

	// add genesis attendees to the app state
	cdc.MustUnmarshalJSON(appState[longy.ModuleName], &genesisState)

	//get the attendees from eventbrite
	genesisState.Attendees, err = utils.GetAttendees()
	if err != nil {
		return genesisState, err
	}

	return genesisState, err
}
