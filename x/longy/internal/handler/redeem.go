package handler

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsgRedeem processes MsgRedeem message in order to set an attendee's winnings as claimed
//nolint:gocritic
func HandleMsgRedeem(ctx sdk.Context, k keeper.Keeper, msg types.MsgRedeem) sdk.Result {
	if !k.IsServiceAccount(ctx, msg.Sender) {
		return types.ErrInsufficientPrivileges("only the service account can call this").Result()
	}

	err := k.RedeemPrizes(ctx, msg.Attendee)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
