package genesis

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
)

// AddSetGenesisKeyServiceCmd will set the testing master public/address keys where "master" is the seed
func AddSetGenesisKeyServiceCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set-genesis-key-service",
		RunE: func(_ *cobra.Command, args []string) error {
			return setGenesisService(ctx, cdc, types.ServiceSeed)
		},
	}

	return cmd
}

// AddSetGenesisBonusServiceCmd -
func AddSetGenesisBonusServiceCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set-genesis-bonus-service",
		RunE: func(_ *cobra.Command, args []string) error {
			return setGenesisService(ctx, cdc, types.BonusServiceSeed)
		},
	}

	return cmd
}

//AddSetGenesisClaimServiceCmd sets the claim service account for prize redemption
func AddSetGenesisClaimServiceCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set-genesis-claim-service",
		RunE: func(_ *cobra.Command, args []string) error {
			return setGenesisService(ctx, cdc, types.ClaimServiceSeed)
		},
	}

	return cmd
}

func setGenesisService(ctx *server.Context, cdc *codec.Codec, seed string) error {
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
	err = cdc.UnmarshalJSON(appState[longy.ModuleName], &genState)
	if err != nil {
		return err
	}
	pubKey := tmcrypto.GenPrivKeySecp256k1([]byte(seed)).PubKey()
	sdkAddr := sdk.AccAddress(pubKey.Address())

	switch seed {
	case types.ServiceSeed:
		genState.KeyService.Address = sdkAddr
		genState.KeyService.PubKey = pubKey
	case types.BonusServiceSeed:
		genState.BonusService.Address = sdkAddr
		genState.BonusService.PubKey = pubKey
	case types.ClaimServiceSeed:
		genState.ClaimService.Address = sdkAddr
		genState.ClaimService.PubKey = pubKey
	default:
		return fmt.Errorf("seed %s is not allowed for a service account", seed)
	}

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
