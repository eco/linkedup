package longy

import (
	"fmt"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler constructor for our button module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) { //cast message

		case types.MsgClaimID:
			return handleMsgClaimID(&ctx, keeper, msg)
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

// handleMsgClaimID processes MsgClaimID in order to associate an address with an id
// nolint: unparam
func handleMsgClaimID(context *sdk.Context, keeper Keeper, id types.MsgClaimID) sdk.Result {
	return sdk.Result{}
}

// handleMsgQrScan processes MsgQrScan
// nolint: unparam
func handleMsgQrScan(ctx *sdk.Context, keeper Keeper, msg types.MsgQrScan) sdk.Result {
	//validate sender address is correct

	//get the scanned address from the QR code

	//update scan state

	return sdk.Result{}
}

// nolint: unparam
func handleMsgShareInfo(ctx *sdk.Context, keeper Keeper, msg types.MsgShareInfo) sdk.Result {

	//update scan state

	return sdk.Result{}
}
