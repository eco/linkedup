package internal

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
		queryType := path[0]
		queryArgs := path[1:]

		switch queryType {
		case AttendeeStoreKey:
			return queryAttendee(&ctx, queryArgs, &keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
		}
	}
}

// nolint
func queryAttendee(ctx *sdk.Context, path []string, keeper *Keeper) (attendee []byte, err sdk.Error) {
	return
}

// nolint
func scanQR(ctx *sdk.Context, path []string, req *abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	return
}
