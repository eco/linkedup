package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryAttendees = "attendees"
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
		default:
			break
		}

		return nil, sdk.ErrUnknownRequest("unknown query endpoint")
	}
}

func queryAttendees(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	//addr, err := sdk.AccAddressFromBech32(path[0])
	//if err != nil {
	//	return nil,  sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos address : %s", path[0]))
	//}
	attendee, ok := keeper.GetAttendeeWithID(ctx, path[0])

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that address")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, attendee)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
