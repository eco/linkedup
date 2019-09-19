package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	appName = "longy"
	//ModuleName us tge bane if this app
	ModuleName = appName
)

// nolint: unused
// base inherits the base app and contains all the KVStore keys as well as the keepers
// to other modules we use.
type base struct {
	// inherits from and it reflects the ABCI application implementation.
	*baseapp.BaseApp
	// the modules that are added by child structs
	modules ModuleInterface
	// the codec to use for encoding/decoding of ABCI messages
	Cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// keeper for managing accounts
	accountKeeper auth.AccountKeeper
	// keeper for managing balances and send/receive of funds
	bankKeeper bank.Keeper
	// keeper for the supply
	supplyKeeper supply.Keeper
	// keeper for managing the params of other modules
	paramsKeeper sdkparams.Keeper

	// Module Manager
	mm *module.Manager
}

// GenericGenesisState represents chain state at the start of the chain.
// Any initial state (account balances) are stored here.
type GenericGenesisState struct {
	// AuthData is the auth state of fees collected at genesis, used by the auth keeper
	AuthData auth.GenesisState `json:"auth"`
	// Accounts is an array of accounts that should exist at genesis
	Accounts []*auth.BaseAccount `json:"accounts"`
	// BandData is an array of accounts and their coin balances at genesis
	BankData bank.GenesisState `json:"bank"`
}

// InitApp is the constructor for out app, it initializes the base app instance, all of the keepers
// for the modules we import as well as our own, sets the routes to our module, and then mounts the stores
// and initializes the chain.
func (app *base) InitApp(modules ModuleInterface, logger log.Logger, db dbm.DB) {
	// First define the top level codec that will be shared by the different modules
	app.modules = modules

	cdc := app.RegisterCodec(nil)

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	// Initialize the application with the store keys it requires
	app.BaseApp = bApp
	app.Cdc = cdc

	app.AddStores()

	app.AddKeepers()

	app.AddModules()

	// The initChain handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChain)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.MountExtraStores()
}

func (app *base) LoadLatest() {
	// loaded after all the stores are mounted
	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.LastBlockHeight()
}

// nolint: gocritic
// initChain unmarshals the genesis state on chain start up, and populates the initial state of the chain. In our
// case it creates the accounts and assigns them the balance that the genesis file has set
func (app *base) initChain(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// unmarshal the genesis state from the file json
	stateJSON := req.AppStateBytes
	genesisState := new(GenericGenesisState)
	err := app.Cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	// add all the accounts to the chain
	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// initialize the modules with the genesis state
	auth.InitGenesis(ctx, app.accountKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)

	return abci.ResponseInitChain{}
}

// nolint: gocritic
func (app *base) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// nolint: gocritic
func (app *base) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// nolint: gocritic
func (app *base) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ExportAppStateAndValidators generates a state file by dumping the current state of the chain
// and all of the accounts and their balances.
func (app *base) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator,
	err error) {
	ctx := app.NewContext(true, abci.Header{})
	var accounts []*auth.BaseAccount

	// function to add accounts
	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	// add all the accounts to our local slice
	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenericGenesisState{
		Accounts: accounts,
		AuthData: auth.DefaultGenesisState(),
		BankData: bank.DefaultGenesisState(),
	}

	// marshal state into a byte array
	appState, err = codec.MarshalJSONIndent(app.Cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}
