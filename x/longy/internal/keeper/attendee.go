package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
)

// GetAttendeeByID will retrieve the attendee by `id`. The Address of an attendee is generated using
// the secp256k1 key using `id` as the secret. returns false if the attendee does not exist
func (k Keeper) GetAttendeeByID(ctx sdk.Context, id string) (types.Attendee, bool) {
	address := util.IDToAddress(id)
	return k.GetAttendeeByAddress(ctx, address)
}

func (k Keeper) GetAttendeeByAddress(ctx sdk.Context, address sdk.AccAddress) (types.Attendee, bool) {
	key := types.AttendeeKey(address)
	bz := k.Get(ctx, key)
	if bz == nil {
		return types.Attendee{}, false
	}

	var a types.Attendee
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &a)
	return a, true
}

// SetAttendee will set the attendee `a` to the store using it's address
func (k Keeper) SetAttendee(ctx sdk.Context, a types.Attendee) {
	addr := a.Address()
	key := types.AttendeeKey(addr)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(a)
	k.Set(ctx, key, bz)
}
