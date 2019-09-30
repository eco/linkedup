package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

// GetAttendee will retrieve the attendee by `id`. returns false if the attendee does not exist
func (k Keeper) GetAttendeeByID(ctx sdk.Context, id string) (types.Attendee, bool) {
	key := types.AttendeeKeyByID(id)
	bz := k.Get(ctx, key)
	if bz == nil {
		return types.Attendee{}, false
	}

	var a types.Attendee
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &a)
	return a, true
}

// SetAttendee will set the attendee `a` by it's id
func (k Keeper) SetAttendee(ctx sdk.Context, a types.Attendee) {
	key := types.AttendeeKeyByID(a.ID)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(a)
	k.Set(ctx, key, bz)
}
