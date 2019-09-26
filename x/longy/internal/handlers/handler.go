package handlers

import (
	"fmt"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler constructor for our button module
func NewHandler(keeper longy.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) { //cast message

		case types.MsgClaimID:
			return handleClaimIDMsg(&ctx, &keeper, &msg)
		case types.MsgQrScan:
			return handleMsgQrScan(&ctx, keeper, msg)
		case types.MsgShareInfo:
			return handleMsgShareInfo(&ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized button Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgQrScan processes MsgQrScan
// nolint: unparam
func handleMsgQrScan(ctx *sdk.Context, keeper longy.Keeper, msg types.MsgQrScan) sdk.Result {
	//validate sender address is correct

	//get the scanned address from the QR code

	//update scan state

	return sdk.Result{}
}

// nolint: unparam
func handleMsgShareInfo(ctx *sdk.Context, keeper longy.Keeper, msg types.MsgShareInfo) sdk.Result {

	//update scan state

	return sdk.Result{}
}
