package client

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// AddGenesisAttendeesCmd returns add-genesis-attendees cobra Command. Allows users to add the list of attendees
// to the chain by their eventbrite id
func AddGenesisAttendeesCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-attendees",
		Short: "Add genesis attendees to genesis.json",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// add genesis attendees to the app state
			var genesisState longy.GenesisState

			cdc.MustUnmarshalJSON(appState[longy.ModuleName], &genesisState)

			genesisState.Attendees, err = utils.GetAttendees()
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
		},
	}

	return cmd
}
