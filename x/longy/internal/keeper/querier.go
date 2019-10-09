package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	//QueryAttendees  is the key for attendee gets
	QueryAttendees = "attendees"
	//QueryScans is the key for scan gets
	QueryScans = "scans"
	//AddressKey is the key for attendee gets by address
	AddressKey = "address"
)

// NewQuerier is the module level router for state queries
//nolint:gocritic
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		queryType := path[0]
		queryArgs := path[1:]

		switch queryType {
		case QueryAttendees:
			if path[1] == AddressKey {
				queryArgs = path[2:]
				return queryAttendeesByAddr(ctx, queryArgs, req, keeper)
			}
			return queryAttendees(ctx, queryArgs, req, keeper)

		case QueryScans:
			return queryScans(ctx, queryArgs, req, keeper)
		default:
			break
		}

		return nil, sdk.ErrUnknownRequest("unknown query endpoint")
	}
}

//nolint:gocritic,unparam
func queryAttendees(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	attendee, ok := keeper.GetAttendeeWithID(ctx, path[0])

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that address")
	}

	res, e := codec.MarshalJSONIndent(keeper.cdc, attendee)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

//nolint:gocritic,unparam
func queryAttendeesByAddr(ctx sdk.Context, path []string, req abci.RequestQuery,
	keeper Keeper) (res []byte, err sdk.Error) {
	addr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos address : %s", path[0]))
	}

	attendee, ok := keeper.GetAttendee(ctx, addr)

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that address")
	}

	res, e = codec.MarshalJSONIndent(keeper.cdc, attendee)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return
}

//nolint:gocritic,unparam
func queryScans(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	scan, err := keeper.GetScanByID(ctx, types.Decode(path[0]))
	if err != nil {
		return
	}
	res, e := codec.MarshalJSONIndent(keeper.cdc, scan)
	if e != nil {
		panic("could not marshal result to JSON")
	}
	return
}
