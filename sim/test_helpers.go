package sim

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// setup initializes a new SimApp. A Nop logger is set in SimApp.
func setup(isCheckTx bool) *LongyApp {
	db := dbm.NewMemDB()
	longyApp := NewLongyApp(log.NewNopLogger(), db, nil, true, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(MakeCodec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		longyApp.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return longyApp
}

// CreateTestApp sets up a new test app
func CreateTestApp(isCheckTx bool) (*LongyApp, sdk.Context) {
	longyApp := setup(isCheckTx)
	ctx := longyApp.BaseApp.NewContext(isCheckTx, abci.Header{})

	return longyApp, ctx
}

/*

// CreateContextAndKeepers will create a in-memory backed keepers and sdk.Context for
// message handler testing
func CreateContextAndKeepers() (sdk.Context, longy.Keeper, auth.AccountKeeper) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NopLogger())
	keys := sdk.NewKVStoreKeys("longy", "auth")

	authSubspace := params.NewSubspace(app.MakeCodec(), params.StoreKey, params.TStoreKey, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(auth.ModuleCdc, keys["auth"], authSupspace, auth.ProtoBaseAccount)
	longyKeeper := longy.NewKeeper(longy.ModuleCdc, keys["longy"], accountKeeper)

	return ctx, longyKeeper, accountKeeper

}
*/
