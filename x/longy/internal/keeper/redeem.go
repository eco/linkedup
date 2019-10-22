package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//RedeemPrizes sets all of the prizes for an attendee to claimed = true
//nolint:gocritic
func (k *Keeper) RedeemPrizes(ctx sdk.Context, attendeeAddr sdk.AccAddress) sdk.Error {
	//get the AccAddress for the scanned qr code
	attendee, ok := k.GetAttendee(ctx, attendeeAddr)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee")
	}

	winnings := attendee.Winnings
	for i := range winnings {
		winnings[i].Claimed = true
	}

	k.SetAttendee(ctx, &attendee)
	return nil
}
