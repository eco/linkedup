package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/onsi/gomega"
)

//AddAttendeeToKeeper is a helper for adding an attendee and its associate account to a test keeper
//nolint:gocritic
func AddAttendeeToKeeper(ctx sdk.Context, keeper *longy.Keeper, badgeID string, claimed bool,
	sponsor bool) (attendee types.Attendee) {
	addr := util.IDToAddress(badgeID)
	acc := keeper.AccountKeeper().NewAccountWithAddress(ctx, addr)
	attendee = types.NewAttendee(badgeID)
	attendee.Sponsor = sponsor
	attendee.Claimed = claimed
	keeper.AccountKeeper().SetAccount(ctx, acc)
	keeper.SetAttendee(ctx, &attendee)
	return
}

//SetRedeemAccount creates and sets an account to be the redeemer
//nolint:gocritic
func SetRedeemAccount(ctx sdk.Context, k longy.Keeper, addresses sdk.AccAddress) exported.Account {
	acc := k.AccountKeeper().NewAccountWithAddress(ctx, addresses)
	k.AccountKeeper().SetAccount(ctx, acc)
	err := k.SetRedeemAccount(ctx, addresses)
	gomega.Expect(err).To(gomega.BeNil())
	return acc
}
