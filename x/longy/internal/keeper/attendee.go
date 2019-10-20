package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	//LeaderBoardCount is the total number of attendees on the board
	LeaderBoardCount = 30
	//LeaderBoardTier1Prize is the tier 1 prize pool in USD
	LeaderBoardTier1Prize = 5000
	//LeaderBoardTier2Prize is the tier 2 prize pool in USD
	LeaderBoardTier2Prize = 2000
	//LeaderBoardTier1Count is the number of people in the first tier
	LeaderBoardTier1Count = 10
)

//Tier is a prize tier in the leader  board
type Tier struct {
	PrizeAmount int              `json:"prizeAmount"`
	Attendees   []types.Attendee `json:"attendees"`
}

//LeaderBoard is the leader board struct
type LeaderBoard struct {
	TotalCount int  `json:"totalCount"`
	Tier1      Tier `json:"tier1"`
	Tier2      Tier `json:"tier2"`
}

//NewLeaderBoard returns an initialized leader board with some constants
func NewLeaderBoard(count int, top []types.Attendee) *LeaderBoard {
	var first, second []types.Attendee
	if len(top) >= LeaderBoardCount {
		first = make([]types.Attendee, LeaderBoardTier1Count)
		second = make([]types.Attendee, LeaderBoardCount-LeaderBoardTier1Count)
		copy(first, top[:LeaderBoardTier1Count])
		copy(second, top[LeaderBoardTier1Count:])
	} else {
		first = top
	}
	return &LeaderBoard{
		TotalCount: count,
		Tier1: Tier{
			PrizeAmount: LeaderBoardTier1Prize,
			Attendees:   first,
		},
		Tier2: Tier{
			PrizeAmount: LeaderBoardTier2Prize,
			Attendees:   second,
		},
	}
}

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

	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &attendee)
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

			err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &attendee)
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

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(a)
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

	multiplier := uint(1)
	bonus := k.GetBonus(ctx)
	if bonus != nil {
		multiplier = bonus.Multiplier
	}

	a1Points := types.ScanAttendeeAwardPoints
	a2Points := types.ScanAttendeeAwardPoints
	if a2.Sponsor {
		a1Points = types.ScanSponsorAwardPoints * multiplier
	}
	if a1.Sponsor {
		a2Points = types.ScanSponsorAwardPoints * multiplier
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
		prize, err := k.GetPrize(ctx, types.GetPrizeIDByTier(attendee.GetTier()))
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

	multiplier := uint(1)
	bonus := k.GetBonus(ctx)
	if bonus != nil {
		multiplier = bonus.Multiplier
	}

	//give sender points for sharing, check if receiver is a sponsor
	val := types.ShareAttendeeAwardPoints
	if receiver.Sponsor {
		val = types.ShareSponsorAwardPoints * multiplier
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
