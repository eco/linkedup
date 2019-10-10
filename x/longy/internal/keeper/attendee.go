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
func (k Keeper) AwardScanPoints(ctx sdk.Context, scan *types.Scan) sdk.Error {
	a1, a2, err := k.getAttendeesByScan(ctx, scan)
	if err != nil {
		return err
	}

	points := types.ScanAttendeeAwardPoints
	if a1.Sponsor || a2.Sponsor {
		points = types.ScanSponsorAwardPoints
	}
	a1.AddRep(points)
	a2.AddRep(points)
	k.SetAttendee(ctx, a1)
	k.SetAttendee(ctx, a2)

	//update scan points
	scan.AddPoints(points, points)
	k.SetScan(ctx, scan)
	return nil
}

//AddSharedID adds the scan id to the scan ids array of both the sender and receiver is they don't contain it yet
//nolint:gocritic
func (k Keeper) AddSharedID(ctx sdk.Context, senderAddr sdk.AccAddress, receiverAddr sdk.AccAddress,
	scanID []byte) sdk.Error {
	sender, receiver, err := k.GetAttendees(ctx, senderAddr, receiverAddr)
	if err != nil {
		return err
	}
	if sender.AddScanID(scanID) {
		k.SetAttendee(ctx, sender)
	}
	if receiver.AddScanID(scanID) {
		k.SetAttendee(ctx, receiver)
	}
	return nil
}

//AwardShareInfoPoints adds points to the sender of the shared info based on if the receiver is a sponsor or not
//nolint:gocritic
func (k Keeper) AwardShareInfoPoints(ctx sdk.Context, scan *types.Scan, senderAddr sdk.AccAddress,
	receiverAddr sdk.AccAddress) sdk.Error {
	sender, receiver, err := k.GetAttendees(ctx, senderAddr, receiverAddr)
	if err != nil {
		return err
	}
	//give sender points for sharing, check if receiver is a sponsor
	val := types.ShareAttendeeAwardPoints
	if receiver.Sponsor {
		val = types.ShareSponsorAwardPoints
	}
	sender.AddRep(val)
	k.SetAttendee(ctx, sender)
	//update scan points
	scan.AddPointsToAccount(sender.Address, val)
	k.SetScan(ctx, scan)
	return nil
}

//GetAttendees returns the attendees for the give account addresses
//nolint:gocritic
func (k Keeper) GetAttendees(ctx sdk.Context, acc1 sdk.AccAddress,
	acc2 sdk.AccAddress) (a1 types.Attendee, a2 types.Attendee, err sdk.Error) {
	var exists bool
	a1, exists = k.GetAttendee(ctx, acc1)
	if !exists {
		err = types.ErrAttendeeNotFound("attendee for points award was not found")
		return
	}
	a2, exists = k.GetAttendee(ctx, acc2)
	if !exists {
		err = types.ErrAttendeeNotFound("attendee for points award was not found")
		return
	}
	return
}

//getAttendeesByScan returns  the attendees for the give scan
//nolint:gocritic
func (k Keeper) getAttendeesByScan(ctx sdk.Context, scan *types.Scan) (a1 types.Attendee,
	a2 types.Attendee, err sdk.Error) {
	return k.GetAttendees(ctx, scan.S1, scan.S2)
}
