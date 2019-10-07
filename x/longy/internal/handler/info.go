package handler

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//HandleMsgInfo processes MsgInfo message
//nolint:gocritic
func HandleMsgInfo(ctx sdk.Context, k keeper.Keeper, msg types.MsgInfo) sdk.Result {
	sender, ok := k.GetAttendee(ctx, msg.Sender)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee that sent info").Result()
	}
	receiver, ok := k.GetAttendee(ctx, msg.Receiver)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee to receive info").Result()
	}

	//check that there is no existing info by this sender
	infoID, err := types.GenInfoID(msg.Sender, msg.Receiver)
	if err != nil {
		return err.Result()
	}
	storedInfo, _ := k.GetInfoByID(ctx, infoID)
	if storedInfo != nil {
		return types.ErrInfoAlreadyExists("sender has already shared info with the receiver").Result()
	}

	//check that there is a complete scan event before accepting the share info
	scanID, err := types.GenScanID(msg.Sender, msg.Receiver)
	if err != nil {
		return err.Result()
	}
	scan, err := k.GetScanByID(ctx, scanID)
	if err != nil {
		return err.Result()
	}
	if !scan.Complete {
		return types.ErrInvalidShareForScan(
			"cannot share info when the corresponding scan is incomplete").Result()
	}

	//create info
	info, err := types.NewInfo(msg.Sender, msg.Receiver, msg.Data)
	if err != nil {
		return err.Result()
	}
	k.SetInfo(ctx, &info)

	//add info to receiver attendee struct
	receiver.InfoIDs = append(receiver.InfoIDs, string(info.ID))
	//give sender points for sharing, check if receiver is a sponsor
	val := types.ShareAttendeeAwardPoints
	if receiver.Sponsor {
		val = types.ShareSponsorAwardPoints
	}
	sender.AddRep(val)
	k.SetAttendee(ctx, sender)
	k.SetAttendee(ctx, receiver)
	return sdk.Result{}
}
