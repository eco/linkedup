package longy

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/eco/longy/x/longy/client/cli"
	"github.com/eco/longy/x/longy/client/rest"
	"github.com/eco/longy/x/longy/internal"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is
type AppModuleBasic struct{}

//Name returns the name of the module
func (a AppModuleBasic) Name() string {
	return ModuleName
}

//RegisterCodec registers module to the codec
func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

//DefaultGenesis returns the default genesis for this module if any
func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

//ValidateGenesis validates that the json genesis is valid to our module
func (a AppModuleBasic) ValidateGenesis(json json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(json, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshaled, passes along to genesis.go
	return ValidateGenesis(data)
}

//RegisterRESTRoutes registers our module rest endpoints
//nolint: gocritic
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, StoreKey)
}

//GetTxCmd returns any tx commands from this module to the parent command in the cli
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

//GetQueryCmd returns any query commands from this module to the parent command in the cli
func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

// AppModule structure holding or keepers together
type AppModule struct {
	AppModuleBasic
	accountKeeper types.AccountKeeper
	bankKeeper    bank.Keeper
	longyKeeper   internal.Keeper
}

// NewAppModule creates a new AppModule object
// nolint: gocritic
func NewAppModule(account types.AccountKeeper, bank bank.Keeper, longy Keeper) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  account,
		bankKeeper:     bank,
		longyKeeper:    longy,
	})
}

// RegisterInvariants registers something
//nolint: gocritic
func (am AppModule) RegisterInvariants(sdk.InvariantRegistry) {
}

// Route returns the route path to this module for transactions
// nolint: gocritic
func (am AppModule) Route() string {
	return RouterKey
}

// NewHandler returns the handler for the module
// nolint: gocritic
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.longyKeeper)
}

// QuerierRoute returns the route path to this module for queries
// nolint: gocritic
func (am AppModule) QuerierRoute() string {
	return ModuleName
}

//NewQuerierHandler returns the handler for queries
// nolint: gocritic
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return internal.NewQuerier(am.longyKeeper)
}

//BeginBlock is the callback for the start of a block
// nolint: gocritic
func (am AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
}

//EndBlock is the callback for the end of a block
// nolint: gocritic
func (am AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// InitGenesis init-genesis
// nolint: gocritic
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis export genesis
// nolint: gocritic
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return make([]byte, 0)
}
