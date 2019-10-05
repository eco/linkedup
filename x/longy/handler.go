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
		case types.MsgRekey:
			return handleMsgRekey(ctx, keeper, msg)
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
func handleMsgRekey(ctx sdk.Context, k Keeper, msg types.MsgRekey) sdk.Result {
	accountKeeper := k.AccountKeeper()

	// authorization passed, we simply need to update the attendee's public key
	acc := accountKeeper.GetAccount(ctx, msg.AttendeeAddress)
	_ = acc.SetPubKey(msg.NewAttendeePublicKey)
	accountKeeper.SetAccount(ctx, acc)

	// update the attendee to unclaimed
	attendee, ok := k.GetAttendee(ctx, msg.AttendeeAddress)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee").Result()
	}
	attendee.SetCommitment(msg.Commitment)
	attendee.SetUnclaimed()
	k.SetAttendee(ctx, attendee)

	return sdk.Result{}
}

//nolint: unparam, gocritic
func handleMsgClaimKey(ctx sdk.Context, k Keeper, msg types.MsgClaimKey) sdk.Result {
	attendee, ok := k.GetAttendee(ctx, msg.AttendeeAddress)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee").Result()
	}

	if attendee.IsClaimed() {
		return types.ErrAttendeeAlreadyClaimed().Result()
	}

	if !attendee.CurrentCommitment().VerifyReveal(msg.Secret) {
		return types.ErrInvalidCommitmentReveal().Result()
	}

	// all checks passed. mark the attendee as claimed
	attendee.ResetCommitment()
	attendee.SetClaimed()
	if !attendee.HasPreviouslyClaimed() {
		// TODO: this probably will change
		attendee.AddRep(uint(5))
	}

	k.SetAttendee(ctx, attendee)

	return sdk.Result{}
}
