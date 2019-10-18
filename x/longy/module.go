package longy

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/eco/longy/x/longy/client/rest"
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

// Name returns the name of the module
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers module to the codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns the default genesis for this module if any
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	gen := DefaultGenesisState()
	return ModuleCdc.MustMarshalJSON(gen)
}

// ValidateGenesis validates that the json genesis is valid to our module
func (AppModuleBasic) ValidateGenesis(data json.RawMessage) error {
	var gen GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &gen)

	return ValidateGenesis(gen)
}

// RegisterRESTRoutes registers our module rest endpoints
//nolint:gocritic
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, StoreKey)
}

// GetTxCmd returns any tx commands from this module to the parent command in the cli
func (AppModuleBasic) GetTxCmd(*codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd returns any query commands from this module to the parent command in the cli
func (AppModuleBasic) GetQueryCmd(*codec.Codec) *cobra.Command {
	return nil
}

// AppModule structure holding or keepers together
type AppModule struct {
	AppModuleBasic

	keeper Keeper
}

// NewAppModule creates a new AppModule object
// nolint: gocritic
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Route returns the route key for the the appropriate messages
//nolint:gocritic
func (AppModule) Route() string {
	return ModuleName
}

// NewHandler return's the module's handler
//nolint:gocritic
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the key for the router
//nolint:gocritic
func (AppModule) QuerierRoute() string {
	return ModuleName
}

// NewQuerierHandler returns the handler for the querier
//nolint:gocritic
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// BeginBlock runs at the beginning of each block
//nolint:gocritic
func (AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
}

// EndBlock runs at the end of each block
//nolint:gocritic
func (AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// RegisterInvariants registers the invariants for this module
//nolint:gocritic
func (AppModule) RegisterInvariants(sdk.InvariantRegistry) {
}

// InitGenesis init-genesis
// nolint: gocritic
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var gen GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &gen)

	InitGenesis(ctx, am.keeper, gen)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis export genesis
// nolint: gocritic
func (AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	//return make([]byte, 0)
	return nil
}
