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
	//QueryBonus is the key for checking for the current bonus
	QueryBonus = "bonus"
	//AddressKey is the key for attendee gets by address
	AddressKey = "address"
	//PrizesKey is the key for the event prizes
	PrizesKey = "prizes"
	//RedeemKey is the key for the redeem event
	RedeemKey = "redeem"
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
				return queryAttendeesByAddr(ctx, queryArgs, keeper)
			}
			return queryAttendees(ctx, queryArgs, keeper)
		case QueryScans:
			return queryScans(ctx, queryArgs, keeper)
		case PrizesKey:
			return queryPrizes(ctx, keeper)
		case QueryBonus:
			return queryBonus(ctx, keeper)
		case RedeemKey:
			return queryRedeem(ctx, keeper, queryArgs)
		default:
			break
		}

		return nil, sdk.ErrUnknownRequest("unknown query endpoint")
	}
}

//nolint:gocritic,unparam
func queryRedeem(ctx sdk.Context, keeper Keeper, path []string) (res []byte, err sdk.Error) {
	addr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos AccAddress : %s", path[0]))
	}
	attendee, ok := keeper.GetAttendee(ctx, addr)
	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that AccAddress")
	}

	//validate sig
	err = ValidateSig(attendee.PubKey, path[0], path[1])
	if err != nil {
		return
	}

	winnings := attendee.Winnings
	ws := make([]types.Win, len(attendee.Winnings))
	for i := range winnings {
		if !winnings[i].Claimed {
			ws = append(ws, winnings[i])
		}
	}

	res, e = codec.MarshalJSONIndent(keeper.cdc, ws)

	if e != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

//nolint:gocritic,unparam
func queryPrizes(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	prizes, err := keeper.GetPrizes(ctx)
	if err != nil {
		return
	}

	res, e := codec.MarshalJSONIndent(keeper.cdc, prizes)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return
}

//nolint:gocritic,unparam
func queryBonus(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	bonus := keeper.GetBonus(ctx)
	if bonus == nil {
		return nil, nil
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, bonus)
	if err != nil {
		panic(fmt.Sprintf("json marshal bonus: %s", err))
	}

	return res, nil
}

//nolint:gocritic,unparam
func queryAttendees(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {

	attendee, ok := keeper.GetAttendeeWithID(ctx, path[0])

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that AccAddress")
	}

	res, e := codec.MarshalJSONIndent(keeper.cdc, attendee)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

//nolint:gocritic,unparam
func queryAttendeesByAddr(ctx sdk.Context, path []string,
	keeper Keeper) (res []byte, err sdk.Error) {
	addr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("cannot turn param into cosmos AccAddress : %s", path[0]))
	}

	attendee, ok := keeper.GetAttendee(ctx, addr)

	if !ok {
		return nil, types.ErrAttendeeNotFound("could not find attendee with that AccAddress")
	}

	res, e = codec.MarshalJSONIndent(keeper.cdc, attendee)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return
}

//nolint:gocritic,unparam
func queryScans(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
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
