package handler

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsgQrScan processes MsgScanQr message
//nolint: unparam, gocritic
func HandleMsgQrScan(ctx sdk.Context, k keeper.Keeper, msg types.MsgScanQr) sdk.Result {
	//get the address for the scanned qr code
	attendee, ok := k.GetAttendeeWithID(ctx, msg.ScannedQR)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee").Result()
	}
	//get the id for the scan event
	id, err := types.GenID(msg.Sender, attendee.Address)
	if err != nil {
		return err.Result()
	}

	//get the scan event
	scan, err := k.GetScanByID(ctx, id)
	if err != nil { //if new scan, create it
		scan, err = types.NewScan(msg.Sender, attendee.Address)
		if err != nil {
			return err.Result()
		}
	} else { //scan already existed
		//since S2 is always the person who's badge was scanned, then both players have scanned
		//if the sender is that person. We can mark off this scan as complete
		if scan.S2.Equals(msg.Sender) && !scan.Complete {
			scan.Complete = true
			//todo reward points?
		} else { //the original scanner must have rescanned the same badge
			return types.ErrScanQRAlreadyOccurred("the scan already exists").Result()
		}
	}

	k.SetScan(ctx, &scan)

	return sdk.Result{}
}
