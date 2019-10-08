package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryAttendees = "attendees"
	QueryScans     = "scans"
)

// NewQuerier is the module level router for state queries
//nolint:gocritic
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		queryType := path[0]
		queryArgs := path[1:]

		switch queryType {
		case QueryAttendees:
			return queryAttendees(ctx, queryArgs, req, keeper)
		case QueryScans:
			return queryScans(ctx, queryArgs, req, keeper)
		default:
			break
		}

		return nil, sdk.ErrUnknownRequest("unknown query endpoint")
	}
}

func queryScans(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	scan, err := keeper.GetScanByID(ctx, []byte(path[0]))
	if err != nil {
		return
	}
	res, e := codec.MarshalJSONIndent(keeper.cdc, scan)
	if e != nil {
		panic("could not marshal result to JSON")
	}
	return
}

func queryAttendees(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	//addr, err := sdk.AccAddressFromBech32(path[0])
	//if err != nil {
	//	return nil,  sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos address : %s", path[0]))
	//}
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
