package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"math"
)

// GetAttendeeWithID will retrieve the attendee by `id`. The Address of an attendee is generated using
// the secp256k1 key using `id` as the secret. returns false if the attendee does not exist
//nolint:gocritic
func (k *Keeper) GetAttendeeWithID(ctx sdk.Context, id string) (types.Attendee, bool) {
	address := util.IDToAddress(id)
	return k.GetAttendee(ctx, address)
}

// GetAttendee will retrieve the attendee via `AccAddress`
//nolint:gocritic
func (k *Keeper) GetAttendee(ctx sdk.Context, address sdk.AccAddress) (attendee types.Attendee, exists bool) {
	key := types.AttendeeKey(address)
	bz, _ := k.Get(ctx, key)
	if bz == nil {
		return
	}

	err := k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &attendee)
	if err != nil {
		panic(err)
	}

	exists = true
	return
}

//GetAllAttendees fetches all the attendees from the kvStore and returns them
//nolint:gocritic
func (k *Keeper) GetAllAttendees(ctx sdk.Context) (attendees []types.Attendee) {
	it := k.KVStore(ctx).Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		key := it.Key()
		if types.IsAttendeeKey(key) {
			var attendee types.Attendee
			bz, _ := k.Get(ctx, key)
			if bz == nil {
				continue
			}

			err := k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &attendee)
			if err != nil {
				continue
			}
			attendees = append(attendees, attendee)
		}
	}

	return
}

// SetAttendee will set the attendee `a` to the store using it's AccAddress
//nolint:gocritic
func (k *Keeper) SetAttendee(ctx sdk.Context, a *types.Attendee) {
	addr := a.GetAddress()
	key := types.AttendeeKey(addr)

	bz, err := k.Cdc.MarshalBinaryLengthPrefixed(a)
	if err != nil {
		panic(err)
	}
	k.Set(ctx, key, bz)
}

//AwardScanPoints awards the points to each participant of the scan
//nolint:gocritic
func (k *Keeper) AwardScanPoints(ctx sdk.Context, scan *types.Scan) sdk.Error {
	a1, a2, err := k.getAttendeesByScan(ctx, scan)
	if err != nil {
		return err
	}

	if !scan.Accepted {
		return types.ErrScanNotAccepted("scan must be accepted by both parties before awarding points")
	}

	multiplier := float64(1)
	bonus := k.GetBonus(ctx)
	if bonus != nil {
		multiplier = bonus.GetMultiplier()
	}

	a1Points := types.ScanAttendeeAwardPoints
	a2Points := types.ScanAttendeeAwardPoints
	if a2.Sponsor && !a1.Sponsor {
		a1Points = uint(math.Floor(float64(types.ScanSponsorAwardPoints) * multiplier))
	}
	if a1.Sponsor && !a2.Sponsor {
		a2Points = uint(math.Floor(float64(types.ScanSponsorAwardPoints) * multiplier))
	}

	err = k.AddRep(ctx, &a1, a1Points)
	if err != nil {
		return err
	}
	err = k.AddRep(ctx, &a2, a2Points)
	if err != nil {
		return err
	}

	//update scan points
	scan.AddPoints(a1Points, a2Points)
	k.SetScan(ctx, scan)
	return nil
}

//AddRep adds reputation to the attendee, and if that pushes them past a tier, then they will be rewarded a prize
//if there are any left for that tier level
//nolint:gocritic
func (k *Keeper) AddRep(ctx sdk.Context, attendee *types.Attendee, points uint) sdk.Error {
	before := attendee.GetTier()
	attendee.AddRep(points)

	if attendee.GetTier() > before {
		for i := before + 1; i <= attendee.GetTier(); i++ {
			prize, err := k.GetPrize(ctx, types.GetPrizeIDByTier(i))
			if err != nil {
				return err
			}
			if prize.Quantity > 0 {
				added := attendee.AddWinning(&types.Win{
					Tier:    prize.Tier,
					Name:    prize.PrizeText,
					Claimed: false,
				})
				if added {
					prize.Quantity--
					k.SetPrize(ctx, &prize)
				}
			}
		}
	}
	k.SetAttendee(ctx, attendee)
	return nil
}

//AwardShareInfoPoints adds points to the sender of the shared info based on if the receiver is a sponsor or not
//nolint:gocritic
func (k *Keeper) AwardShareInfoPoints(ctx sdk.Context, scan *types.Scan, senderAddr sdk.AccAddress,
	receiverAddr sdk.AccAddress) sdk.Error {
	sender, receiver, err := k.GetAttendees(ctx, senderAddr, receiverAddr)
	if err != nil {
		return err
	}

	if !scan.Accepted {
		return types.ErrScanNotAccepted("scan must be accepted by both parties before awarding points")
	}

	multiplier := float64(1)
	bonus := k.GetBonus(ctx)
	if bonus != nil {
		multiplier = bonus.GetMultiplier()
	}

	//give sender points for sharing, check if receiver is a sponsor
	val := types.ShareAttendeeAwardPoints
	if receiver.Sponsor && !sender.Sponsor {
		val = uint(math.Floor(float64(types.ShareSponsorAwardPoints) * multiplier))
	}

	err = k.AddRep(ctx, &sender, val)
	if err != nil {
		return err
	}
	//update scan points
	scan.AddPointsToAccount(sender.Address, val)
	k.SetScan(ctx, scan)
	return nil
}

//AddSharedID adds the scan id to the scan ids array of both the sender and receiver is they don't contain it yet
//nolint:gocritic
func (k *Keeper) AddSharedID(ctx sdk.Context, senderAddr sdk.AccAddress, receiverAddr sdk.AccAddress,
	scanID []byte) sdk.Error {
	sender, receiver, err := k.GetAttendees(ctx, senderAddr, receiverAddr)
	if err != nil {
		return err
	}
	if sender.AddScanID(scanID) {
		k.SetAttendee(ctx, &sender)
	}
	if receiver.AddScanID(scanID) {
		k.SetAttendee(ctx, &receiver)
	}
	return nil
}

//GetAttendees returns the attendees for the give account addresses
//nolint:gocritic
func (k *Keeper) GetAttendees(ctx sdk.Context, acc1 sdk.AccAddress,
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
func (k *Keeper) getAttendeesByScan(ctx sdk.Context, scan *types.Scan) (a1 types.Attendee,
	a2 types.Attendee, err sdk.Error) {
	return k.GetAttendees(ctx, scan.S1, scan.S2)
}
