package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
)

//ModuleInterface is the interface we use to make module integration into the base app simple
type ModuleInterface interface {
	RegisterCodec(cdc *codec.Codec) *codec.Codec
	AddStores()
	AddKeepers()
	AddRoutes(rtr baseapp.Router)
	AddQueryRoutes(router baseapp.QueryRouter)
	MountExtraStores()
}

func (app *base) AddStores() {
	app.keyMainStore = sdk.NewKVStoreKey("main")
	app.keyAccountStore = sdk.NewKVStoreKey(auth.StoreKey)
	app.keyFeeCollectionStore = sdk.NewKVStoreKey("fee_collection")
	app.keyParamsStore = sdk.NewKVStoreKey("params")
	app.tkeyParamsStore = sdk.NewTransientStoreKey("transient_params")
	if app.modules != nil {
		app.modules.AddStores()
	}
}

func (app *base) AddKeepers() {
	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = sdkparams.NewKeeper(app.Cdc, app.keyParamsStore, app.tkeyParamsStore)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.Cdc,
		app.keyAccountStore,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	// The FeeCollectionKeeper collects transaction fees and renders them to the fee distribution module
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.Cdc, app.keyFeeCollectionStore)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	if app.modules != nil {
		app.modules.AddKeepers()
	}
}

func (app *base) AddRoutes(rtr baseapp.Router) {
	// The app.Router is the main transaction router where each module registers its routes
	// Register the bank and button routes here
	rtr.
		AddRoute("bank", bank.NewHandler(app.bankKeeper))

	if app.modules != nil {
		app.modules.AddRoutes(rtr)
	}
}

func (app *base) AddQueryRoutes(rtr baseapp.QueryRouter) {
	rtr.AddRoute(auth.StoreKey, auth.NewQuerier(app.accountKeeper))
	if app.modules != nil {
		app.modules.AddQueryRoutes(rtr)
	}
}

func (app *base) MountExtraStores() {
	//Mount all the stores to the keys
	app.MountStores(
		app.keyMainStore,
		app.keyAccountStore,
		app.keyFeeCollectionStore,
		app.keyParamsStore,
		app.tkeyParamsStore,
	)

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

func BaseCodec(cdc *codec.Codec) *codec.Codec {
	//accounts and coins
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	//generic tx and crypto messaging
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
