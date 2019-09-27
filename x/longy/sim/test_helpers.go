package sim

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	app "github.com/eco/longy"
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
		stateBytes, err := codec.MarshalJSONIndent(app.MakeCodec(), genesisState)
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

//CreateTestApp sets up a new test app
func CreateTestApp(isCheckTx bool) (*LongyApp, sdk.Context) {
	longyApp := setup(isCheckTx)
	ctx := longyApp.BaseApp.NewContext(isCheckTx, abci.Header{})

	return longyApp, ctx
}
