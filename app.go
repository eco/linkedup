package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"os"
)

const appName = longy.ModuleName

var (
	// DefaultCLIHome is the default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.lycli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.lyd")
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	longy.RegisterCodec(cdc)

	return cdc
}

// LongyApp is our app structure
type LongyApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	longyKeeper longy.Keeper
}

// NewLongyApp is a constructor function for LongyApp
func NewLongyApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *LongyApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, nil, baseAppOptions...)

	longyStoreKey := sdk.NewKVStoreKey("longy")

	var app = &LongyApp{
		BaseApp: bApp,

		longyKeeper: longy.NewKeeper(longyStoreKey, cdc),
		cdc:         cdc,
	}

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(nil)

	// initialize stores
	app.MountStores(longyStoreKey)
	err := app.LoadLatestVersion(longyStoreKey)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}
