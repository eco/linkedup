package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//SetRedeemAccount sets the redeem account from the genesis file
//nolint:gocritic
func (k Keeper) SetRedeemAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	if addr.Empty() {
		return sdk.ErrInvalidAddress(addr.String())
	}

	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	key := types.RedeemKey()
	bz := addr.Bytes()
	k.Set(ctx, key, bz)

	return nil
}

//IsRedeemAccount returns true if the the account passed in is the redeemer account
//nolint:gocritic
func (k Keeper) IsRedeemAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.RedeemKey()
	bz, err := k.Get(ctx, key)
	if err != nil {
		return false
	}
	redeemer := sdk.AccAddress(bz)
	return redeemer.Equals(addr)
}

//RedeemPrizes sets all of the prizes for an attendee to claimed = true
//nolint:gocritic
func (k *Keeper) RedeemPrizes(ctx sdk.Context, attendeeAddr sdk.AccAddress) sdk.Error {
	//get the address for the scanned qr code
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
