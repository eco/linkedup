package querier

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

//nolint:gocritic,unparam
func queryAttendees(ctx sdk.Context, path []string, keeper keeper.Keeper) (res []byte, err sdk.Error) {

	attendee, ok := keeper.GetAttendeeWithID(ctx, path[0])

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that AccAddress")
	}

	res, e := codec.MarshalJSONIndent(keeper.Cdc, attendee)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

//nolint:gocritic,unparam
func queryAttendeesByAddr(ctx sdk.Context, path []string,
	keeper keeper.Keeper) (res []byte, err sdk.Error) {
	addr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos AccAddress : %s", path[0]))
	}

	attendee, ok := keeper.GetAttendee(ctx, addr)

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that AccAddress")
	}

	res, e = codec.MarshalJSONIndent(keeper.Cdc, attendee)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return
}
