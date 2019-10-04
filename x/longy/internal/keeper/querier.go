package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
//nolint:gocritic
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		queryType := path[0]
		//queryArgs := path[1:]

		switch queryType {
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
		}
	}
}
