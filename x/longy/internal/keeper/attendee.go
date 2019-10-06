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
func (k Keeper) GetAttendee(ctx sdk.Context, address sdk.AccAddress) (attendee types.Attendee, exists bool) {
	key := types.AttendeeKey(address)
	bz, _ := k.Get(ctx, key)
	if bz == nil {
		return
	}

	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &attendee)
	if err != nil {
		panic(err)
	}
	exists = true
	return
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

//AwardScanPoints awards the points to each participant of the scan
//nolint:gocritic
func (k Keeper) AwardScanPoints(ctx sdk.Context, scan types.Scan) sdk.Error {
	if !scan.Complete {
		return types.ErrScanNotComplete("cannot reward points for a scan that is not complete")
	}
	a1, exists := k.GetAttendee(ctx, scan.S1)
	if !exists {
		return types.ErrAttendeeNotFound("attendee for points award was not found")
	}
	a2, exists := k.GetAttendee(ctx, scan.S2)
	if !exists {
		return types.ErrAttendeeNotFound("attendee for points award was not found")
	}

	points := types.ScanAttendeeAwardPoints
	if a1.Sponsor || a2.Sponsor {
		points = types.ScanSponsorAwardPoints
	}
	a1.Rep += points
	a2.Rep += points
	k.SetAttendee(ctx, a1)
	k.SetAttendee(ctx, a2)
	return nil
}
