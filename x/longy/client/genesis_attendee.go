package client

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// AddGenesisAttendeesCmd returns add-genesis-attendees cobra Command. Allows users to add the list of attendees
// to the chain by their eventbrite id
func AddGenesisAttendeesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-attendees [serviceAddr]",
		Short: "Add genesis attendees and the service address to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return AddGenesisAttendees(args, ctx, cdc)
		},
	}

	return cmd
}

//AddGenesisAttendees adds the attendees and the service account to the genesis file under the longy key
func AddGenesisAttendees(args []string, ctx *server.Context, cdc *codec.Codec) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))

	// retrieve the app state
	genFile := config.GenesisFile()
	appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
	if err != nil {
		return err
	}

	genesisState, err := BuildGenesisState(appState, cdc, args)
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

//BuildGenesisState builds the genesis state for the longy module
func BuildGenesisState(appState map[string]json.RawMessage, cdc *codec.Codec,
	args []string) (genesisState longy.GenesisState, err sdk.Error) {
	// add genesis attendees to the app state
	cdc.MustUnmarshalJSON(appState[longy.ModuleName], &genesisState)

	//get the service account from the args
	addr, err := getServiceAcct(appState, cdc, args)
	if err != nil {
		return
	}
	genesisState.Service = types.GenesisService{Address: addr}

	//get the attendees from eventbrite
	genesisState.Attendees, err = utils.GetAttendees()
	if err != nil {
		return
	}

	return
}

//getServiceAcct gets the service account from the cmd args and checks that the same account exists in
//the genesis accounts
func getServiceAcct(appState map[string]json.RawMessage, cdc *codec.Codec, args []string) (addr sdk.AccAddress,
	e sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		e = types.ErrGenesisServiceAccountInvalid(err.Error())
		return
	}

	var accountGenesisState genaccounts.GenesisState
	cdc.MustUnmarshalJSON(appState[genaccounts.ModuleName], &accountGenesisState)
	for i := range accountGenesisState {
		if accountGenesisState[i].Address.Equals(addr) {
			return
		}
	}
	e = types.ErrGenesisServiceAccountNotPresent("service account was not found in genesis accounts, " +
		"have you added by calling 'lyd add-genesis-account {acct} {coints...}'")
	return
}
