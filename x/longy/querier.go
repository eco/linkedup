package button

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by this Querier
const (
	QueryPrize = "prize"
	QueryAge   = "age"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		return
	}

}

// nolint: unparam
// resolve: This takes a name and returns the value that is stored by the button. This is similar to a DNS query.
func scanQR(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	return
}
