package querier

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

//nolint:gocritic,unparam
func queryRedeem(ctx sdk.Context, k keeper.Keeper, path []string) (res []byte, err sdk.Error) {
	addr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos AccAddress : %s", path[0]))
	}
	attendee, ok := k.GetAttendee(ctx, addr)
	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that AccAddress")
	}

	////validate sig
	//err = keeper.ValidateSig(attendee.PubKey, path[0], path[1])
	//if err != nil {
	//	return
	//}

	winnings := attendee.Winnings
	ws := make([]types.Win, 0, len(attendee.Winnings))
	for i := range winnings {
		if !winnings[i].Claimed {
			ws = append(ws, winnings[i])
		}
	}

	res, e = codec.MarshalJSONIndent(k.Cdc, ws)

	if e != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
