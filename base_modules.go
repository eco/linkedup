package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

//ModuleInterface is the interface we use to make module integration into the base app simple
type ModuleInterface interface {
	RegisterCodec(cdc *codec.Codec) *codec.Codec
	AddStores()
	AddKeepers()
	GetAppModules() []module.AppModule
	MountExtraStores()
}

func (app *base) AddStores() {
	app.keys = sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey,
		supply.StoreKey, params.StoreKey)
	app.tkeys = sdk.NewTransientStoreKeys(params.TStoreKey)

	if app.modules != nil {
		app.modules.AddStores()
	}
}

func (app *base) AddKeepers() {
	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = sdkparams.NewKeeper(app.Cdc, app.keys[params.StoreKey],
		app.tkeys[params.TStoreKey], sdkparams.DefaultCodespace)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.Cdc,
		app.keys[auth.StoreKey],
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		nil,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	if app.modules != nil {
		app.modules.AddKeepers()
	}
}
func (app *base) AddModules() {
	modules := []module.AppModule{
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
	}

	if app.modules != nil {
		modules = append(modules, app.modules.GetAppModules()...)
	}

	app.mm = module.NewManager(modules...)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	var modNames = make([]string, len(modules))
	for i, mod := range modules {
		modNames[i] = mod.Name()
	}
	app.mm.SetOrderInitGenesis(modNames...)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
}

func (app *base) MountExtraStores() {
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)

	if app.modules != nil {
		app.modules.MountExtraStores()
	}
}

// MakeCodec generates the necessary codecs for Amino
func (app *base) RegisterCodec(cdc *codec.Codec) *codec.Codec {
	//our game logic and state
	if cdc == nil {
		cdc = codec.New()
	}

	cdc = BaseCodec(cdc)

	if app.modules != nil {
		cdc = app.modules.RegisterCodec(cdc)
	}
	return cdc
}

//BaseCodec returns the basic codec for the pre-wired modules
func BaseCodec(cdc *codec.Codec) *codec.Codec {
	//accounts and coins
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	//generic tx and crypto messaging
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
