package genesis

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/eco/longy/x/longy"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
)

func getGenesisState(ctx *server.Context, cdc *codec.Codec) (genesisState map[string]json.RawMessage,
	genDoc *tmtypes.GenesisDoc, genFile string, err error) {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))

	// retrieve the app state
	genFile = config.GenesisFile()
	genesisState, genDoc, err = genutil.GenesisStateFromGenFile(cdc, genFile)
	return
}

//nolint:gocritic
func updateGenesisState(cdc *codec.Codec, genesisState longy.GenesisState, appState map[string]json.RawMessage,
	genDoc *tmtypes.GenesisDoc, genFile string) error {

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
