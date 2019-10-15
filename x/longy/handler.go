package longy

import (
	"fmt"
	"github.com/eco/longy/x/longy/internal/handler"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler constructor for our button module
//nolint:gocritic
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgScanQr:
			return handler.HandleMsgQrScan(ctx, keeper, msg)
		case types.MsgInfo:
			return handler.HandleMsgInfo(ctx, keeper, msg)
		case types.MsgKey:
			return handleMsgKey(ctx, keeper, msg)
		case types.MsgClaimKey:
			return handleMsgClaimKey(ctx, keeper, msg)
		case types.MsgBonus:
			return handleBonus(ctx, keeper, msg)
		case types.MsgClearBonus:
			return handleClearBonus(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s msg type: %T", RouterKey, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
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

	// verify that master key used to sign this message
	masterAddress := k.GetMasterAddress(ctx)
	if masterAddress.Empty() {
		return sdk.ErrInternal("master address not set").Result()
	} else if !masterAddress.Equals(msg.MasterAddress) {
		return sdk.ErrUnauthorized("incorrect master address that signed this message").Result()
	}

	// Check that a public key has not already been set. The rekey service should only be able to
	// submit and alter the public key once
	if account.GetPubKey() != nil {
		return types.ErrAttendeeKeyed("attendee already key'd their account").Result()
	}

	// authorization passed, we simply need to update the attendee's public key
	_ = account.SetPubKey(msg.NewAttendeePublicKey)
	accountKeeper.SetAccount(ctx, account)

	// update the commitment so that the attendee must claim against
	attendee.SetCommitment(msg.Commitment)
	k.SetAttendee(ctx, &attendee)

	return sdk.Result{}
}

//nolint: unparam, gocritic
//pull out and fix tests, test for name, and time stamp
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

	// award rep for the onboarding flow
	err := k.AddRep(ctx, &attendee, types.ClaimBadgeAwardPoints)
	if err != nil {
		return err.Result()
	}

	// add the rsa public key
	attendee.Name = msg.Name
	//Set the time TODO check that this is indeed deterministic time on block header
	attendee.UnixTimeSecClaimed = ctx.BlockTime().Unix()
	attendee.RsaPublicKey = msg.RsaPublicKey
	attendee.EncryptedInfo = msg.EncryptedInfo

	// mark the attendee as claimed
	attendee.SetClaimed()
	k.SetAttendee(ctx, &attendee)

	return sdk.Result{}
}

func handleBonus(ctx sdk.Context, k Keeper, msg types.MsgBonus) sdk.Result {
	// verify that only the master account can send this message
	masterAddress := k.GetMasterAddress(ctx)
	if masterAddress.Empty() {
		return sdk.ErrInternal("master address not set").Result()
	} else if !masterAddress.Equals(msg.MasterAddress) {
		return sdk.ErrUnauthorized("incorrect master address").Result()
	}

	// check if a bonus is already live
	if k.HasLiveBonus(ctx) {
		return sdk.ErrUnauthorized("current bonus must be cleared before setting a new one").Result()
	}

	// set the current bonus
	bonus := types.NewBonus(msg.Multiplier)
	k.SetBonus(ctx, bonus)

	return sdk.Result{}
}

func handleClearBonus(ctx sdk.Context, k Keeper, msg types.MsgClearBonus) sdk.Result {
	// verify that only the master account can send this message
	masterAddress := k.GetMasterAddress(ctx)
	if masterAddress.Empty() {
		return sdk.ErrInternal("master address not set").Result()
	} else if !masterAddress.Equals(msg.MasterAddress) {
		return sdk.ErrUnauthorized("incorrect master address").Result()
	}

	if !k.HasLiveBonus(ctx) {
		return types.ErrDefault("no bonus to clear").Result()
	}

	// clear the bonus
	k.ClearBonus(ctx)

	return sdk.Result{}
}
