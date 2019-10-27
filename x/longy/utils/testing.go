package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/onsi/gomega"
	crypto "github.com/tendermint/tendermint/crypto/secp256k1"
)

//AddAttendeeToKeeper is a helper for adding an attendee and its associate account to a test keeper
//nolint:gocritic
func AddAttendeeToKeeper(ctx sdk.Context, keeper *longy.Keeper, badgeID string, claimed bool,
	sponsor bool) (attendee types.Attendee) {
	addr := util.IDToAddress(badgeID)
	acc := keeper.AccountKeeper().NewAccountWithAddress(ctx, addr)
	attendee = types.NewAttendee(badgeID, sponsor)
	attendee.Claimed = claimed
	keeper.AccountKeeper().SetAccount(ctx, acc)
	keeper.SetAttendee(ctx, &attendee)
	return
}

//SetServiceAccount creates and sets an account to be the service
//nolint:gocritic
func SetServiceAccount(ctx sdk.Context, k longy.Keeper, addresses sdk.AccAddress) exported.Account {
	acc := k.AccountKeeper().NewAccountWithAddress(ctx, addresses)
	pubKey := crypto.GenPrivKeySecp256k1([]byte("service")).PubKey()
	err := acc.SetPubKey(pubKey)
	gomega.Expect(err).To(gomega.BeNil())
	k.AccountKeeper().SetAccount(ctx, acc)
	err = k.SetServiceAddress(ctx, addresses)
	gomega.Expect(err).To(gomega.BeNil())
	return acc
}
