package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
)

// GetAttendeeWithID will retrieve the attendee by `id`. The Address of an attendee is generated using
// the secp256k1 key using `id` as the secret. returns false if the attendee does not exist
//nolint:gocritic
func (k Keeper) GetAttendeeWithID(ctx sdk.Context, id string) (types.Attendee, bool) {
	address := util.IDToAddress(id)
	return k.GetAttendee(ctx, address)
}

// GetAttendee will retrieve the attendee via `address`
//nolint:gocritic
func (k Keeper) GetAttendee(ctx sdk.Context, address sdk.AccAddress) (types.Attendee, bool) {
	key := types.AttendeeKey(address)
	bz := k.Get(ctx, key)
	if bz == nil {
		return types.Attendee{}, false
	}

	var a types.Attendee
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &a)
	if err != nil {
		panic(err)
	}
	return a, true
}

// SetAttendee will set the attendee `a` to the store using it's address
//nolint:gocritic
func (k Keeper) SetAttendee(ctx sdk.Context, a types.Attendee) {
	addr := a.GetAddress()
	key := types.AttendeeKey(addr)

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(a)
	if err != nil {
		panic(err)
	}
	k.Set(ctx, key, bz)
}
