package handler

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsgQrScan processes MsgScanQr message
//nolint:gocritic
func HandleMsgQrScan(ctx sdk.Context, k keeper.Keeper, msg types.MsgScanQr) sdk.Result {
	//get the address for the scanned qr code
	attendee, ok := k.GetAttendeeWithID(ctx, msg.ScannedQR)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee").Result()
	}
	//get the id for the scan event
	id, err := types.GenScanID(msg.Sender, attendee.Address)
	if err != nil {
		return err.Result()
	}

	//get the scan event
	scan, err := k.GetScanByID(ctx, id)
	if err != nil { //if new scan, create it
		err = handleNewScan(ctx, k, msg, attendee)
		if err != nil {
			return err.Result()
		}
	}

	if len(msg.Data) > 0 && scan.{
		err := k.AwardShareInfoPoints(ctx, msg.Sender, attendee.Address)
		if err != nil {
			return err
		}

		err = k.AddSharedInfo(ctx, msg.Sender, attendee.Address, msg.Data)
		if err != nil {
			return err
		}
	}

	//check who is in what position
	var oldData []byte
	if scan.S1.Equals(msg.Sender) {
		oldData = scan.D1
	} else {
		oldData = scan.D2
	}
	if len(oldData) == 0 && len(msg.Data) > 0 {
		err = k.AwardShareInfoPoints(ctx, scan, msg.Sender)
		if err != nil {
			return err
		}
	}

	return sdk.Result{}
}

func handleNewScan(ctx sdk.Context, k keeper.Keeper, msg types.MsgScanQr, attendee types.Attendee) sdk.Error {
	scan, err := types.NewScan(msg.Sender, attendee.Address, msg.Data, nil)

	if err != nil {
		return err
	}
	err = k.AwardScanPoints(ctx, scan)
	if err != nil {
		return err
	}

	k.SetScan(ctx, &scan)
	return nil
}
