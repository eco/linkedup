package longy

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Keeper is temporary
type Keeper struct{}

// NewHandler constructor for our button module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) { //cast message

		case MsgQrScan:
			return handleMsgQrScan(&ctx, keeper, msg)
		case MsgShareInfo:
			return handleMsgShareInfo(&ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized button Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgQrScan processes MsgQrScan
// nolint: unparam
func handleMsgQrScan(ctx *sdk.Context, keeper Keeper, msg MsgQrScan) sdk.Result {
	//validate sender address is correct

	//get the scanned address from the QR code

	//update scan state

	return sdk.Result{}
}

// nolint: unparam
func handleMsgShareInfo(ctx *sdk.Context, keeper Keeper, msg MsgShareInfo) sdk.Result {

	//update scan state

	return sdk.Result{}
}
