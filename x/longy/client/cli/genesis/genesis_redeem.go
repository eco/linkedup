package genesis

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy"
	"github.com/spf13/cobra"
)

// AddGenesisRedeemCmd returns add-redeem-account cobra Command. Allows users to add the redeem account
// that will be used to mark prizes claimed
func AddGenesisRedeemCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-redeem-account [address]",
		Short: "Add redeem account to the genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			return addGenesisRedeemAccount(ctx, cdc, addr)
		},
	}

	return cmd
}

func addGenesisRedeemAccount(ctx *server.Context, cdc *codec.Codec, addr sdk.AccAddress) error {
	appState, genDoc, genFile, err := getGenesisState(ctx, cdc)
	if err != nil {
		return err
	}
	genesisState := buildRedeemGenesisState(appState, cdc, addr)

	return updateGenesisState(cdc, genesisState, appState, genDoc, genFile)
}

func buildRedeemGenesisState(appState map[string]json.RawMessage, cdc *codec.Codec,
	addr sdk.AccAddress) longy.GenesisState {
	var (
		genesisState longy.GenesisState
	)
	fmt.Printf("adding redeem account : %s\n", addr.String())
	// un-marshal the current state of the genesis object
	cdc.MustUnmarshalJSON(appState[longy.ModuleName], &genesisState)

	//get the prizes
	//genesisState.Redeem = longy.GenesisRedeemKey{
	//	Address: addr,
	//}

	return genesisState
}
