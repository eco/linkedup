package handler

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//HandleMsgInfo processes MsgInfo message
//nolint:gocritic
func HandleMsgInfo(ctx sdk.Context, k keeper.Keeper, msg types.MsgInfo) sdk.Result {
	_, receiver, err := k.GetAttendees(ctx, msg.Sender, msg.Receiver)
	if err != nil {
		return err.Result()
	}

	//check that there is an existing scan between these participants
	scanID, err := types.GenScanID(msg.Sender, msg.Receiver)
	if err != nil {
		return err.Result()
	}
	scan, err := k.GetScanByID(ctx, scanID)
	if err != nil {
		return types.ErrScanNotFound("could not find scan for info share").Result()
	}

	err = handleShareInfo(ctx, k, scan, msg.Sender, receiver, msg.Data)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
