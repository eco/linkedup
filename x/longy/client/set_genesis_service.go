package client

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/eco/longy/x/longy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
)

// AddSetGenesisServiceCmd will set the tesing master public/address keys where "master" is the seed
func AddSetGenesisServiceCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set-genesis-service",
		RunE: func(_ *cobra.Command, args []string) error {
			return setGenesisService(ctx, cdc)
		},
	}

	return cmd
}

func setGenesisService(ctx *server.Context, cdc *codec.Codec) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))

	// retrieve genesis
	// retrieve the app state
	genFile := config.GenesisFile()
	appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
	if err != nil {
		return err
	}

	var genState longy.GenesisState
	if err := cdc.UnmarshalJSON(appState[longy.ModuleName], &genState); err != nil {
		return err
	}

	pubKey := tmcrypto.GenPrivKeySecp256k1([]byte("master")).PubKey()
	sdkAddr := sdk.AccAddress(pubKey.Address())
	genState.KeyService.Address = sdkAddr
	genState.KeyService.PubKey = pubKey

	bz, err := cdc.MarshalJSON(genState)
	if err != nil {
		return err
	}
	appState[longy.ModuleName] = bz

	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}
