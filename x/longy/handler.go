package longy

import (
	"fmt"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler constructor for our button module
//nolint:gocritic
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) { //cast message

		case types.MsgQrScan:
			return handleMsgQrScan(ctx, keeper, msg)
		case types.MsgShareInfo:
			return handleMsgShareInfo(ctx, keeper, msg)
		case types.MsgKey:
			return handleMsgKey(ctx, keeper, msg)
		case types.MsgClaimKey:
			return handleMsgClaimKey(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized longy msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgQrScan processes MsgQrScan
//nolint: unparam, gocritic
func handleMsgQrScan(ctx sdk.Context, k Keeper, msg types.MsgQrScan) sdk.Result {
	//validate sender address is correct

	//get the scanned address from the QR code

	//update scan state

	return sdk.Result{}
}

//nolint: unparam, gocritic
func handleMsgShareInfo(ctx sdk.Context, k Keeper, msg types.MsgShareInfo) sdk.Result {

	//update scan state

	return sdk.Result{}
}

//nolint: unparam, gocritic
func handleMsgKey(ctx sdk.Context, k Keeper, msg types.MsgKey) sdk.Result {
	/**
	* For every attendee, there is a cosmos account. This assumption is ensured on `InitGenesis`
	* The following code has checks against both the cosmos account and attendee
	*
	* i.e account == nil || !ok (attendee and the cosmos account exists)
	 */

	// retrieve account/attendee from the store
	accountKeeper := k.AccountKeeper()
	account := accountKeeper.GetAccount(ctx, msg.AttendeeAddress)
	attendee, ok := k.GetAttendee(ctx, msg.AttendeeAddress)
	if account == nil || !ok {
		return types.ErrAttendeeNotFound("nonexistent attendee").Result()
	}

	// Check that a public key has not already been set. The rekey service should only be able to
	// submit and alter the public key once
	if len(account.GetPubKey().Bytes()) > 0 {
		return types.ErrAccountKeyed("attendee already key'd their account").Result()
	}

	// authorization passed, we simply need to update the attendee's public key
	_ = account.SetPubKey(msg.NewAttendeePublicKey)
	accountKeeper.SetAccount(ctx, account)

	// update the commitment so that the attendee must claim against
	attendee.SetCommitment(msg.Commitment)
	k.SetAttendee(ctx, attendee)

	return sdk.Result{}
}

//nolint: unparam, gocritic
func handleMsgClaimKey(ctx sdk.Context, k Keeper, msg types.MsgClaimKey) sdk.Result {

	// retrieve the attendee and make sure the attendee has not been claimed
	attendee, ok := k.GetAttendee(ctx, msg.AttendeeAddress)
	if !ok {
		return types.ErrAttendeeNotFound("nonexistent attendee").Result()
	} else if attendee.IsClaimed() {
		return types.ErrAttendeeClaimed("claimed attendee").Result()
	}

	// verify the commitment
	if !attendee.CurrentCommitment().VerifyReveal(msg.Secret) {
		return types.ErrInvalidCommitmentReveal("incorrect commitment").Result()
	}

	// TODO: disperse reward for onboarding here

	// mark the attendee as claimed
	attendee.SetClaimed()
	k.SetAttendee(ctx, attendee)

	return sdk.Result{}
}
