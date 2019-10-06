package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
)

//AddAttendeeToKeeper is a helper for adding an attendee and its associate account to a test keeper
//nolint:gocritic
func AddAttendeeToKeeper(ctx sdk.Context, keeper *longy.Keeper, badgeID string,
	sponsor bool) (attendee types.Attendee) {
	addr := util.IDToAddress(badgeID)
	acc := keeper.AccountKeeper().NewAccountWithAddress(ctx, addr)
	attendee = types.NewAttendee(badgeID)
	attendee.Sponsor = sponsor
	keeper.AccountKeeper().SetAccount(ctx, acc)
	keeper.SetAttendee(ctx, attendee)
	return
}
