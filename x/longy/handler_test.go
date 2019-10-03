package longy_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/sim"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/stretchr/testify/require"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
	"testing"
)

// will create a new attendee/account with a badge id of "1".
// the corresponding address is returned
func setup(ctx sdk.Context, app *sim.LongyApp) sdk.AccAddress {
	id := "1"
	addr := util.IDToAddress(id)
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	attendee := types.NewAttendee(id)

	app.AccountKeeper.SetAccount(ctx, acc)
	app.LongyKeeper.SetAttendee(ctx, attendee)

	return addr
}

func TestInvalidMsg(t *testing.T) {
	app, ctx := sim.CreateTestApp(true)

	h := longy.NewHandler(app.LongyKeeper)
	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized longy msg type"))
}

func TestRekeyMsg(t *testing.T) {
	app, ctx := sim.CreateTestApp(true)
	addr := setup(ctx, app)

	_, commitment := util.CreateCommitment()
	newPub := tmcrypto.GenPrivKey().PubKey()

	/** setup a rekey against the account **/
	msg := types.MsgRekey{
		AttendeeAddress:      addr,
		NewAttendeePublicKey: newPub,
		Commitment:           commitment,
	}

	h := longy.NewHandler(app.LongyKeeper)
	res := h(ctx, msg)
	require.True(t, res.IsOK())

	// Account should swap keys
	acc := app.AccountKeeper.GetAccount(ctx, addr)
	require.NotNil(t, acc)
	require.Equal(t, newPub.Address(), acc.GetPubKey().Address())

	// Attendee must not be in the claimed state
	a, ok := app.LongyKeeper.GetAttendee(ctx, addr)
	require.True(t, ok)
	require.False(t, a.IsClaimed())
	require.Equal(t, commitment, a.CurrentCommitment())
}

func TestClaimMsg(t *testing.T) {
	app, ctx := sim.CreateTestApp(true)
	addr := setup(ctx, app)

	/** setup a rekey against the account **/
	secret, commitment := util.CreateCommitment()
	newPub := tmcrypto.GenPrivKey().PubKey()

	rekeyMsg := types.MsgRekey{
		AttendeeAddress:      addr,
		NewAttendeePublicKey: newPub,
		Commitment:           commitment,
	}

	h := longy.NewHandler(app.LongyKeeper)
	res := h(ctx, rekeyMsg)
	require.True(t, res.IsOK())

	/** claim the account **/

	// 1. non-existent attendee
	msg := types.MsgClaimKey{
		AttendeeAddress: sdk.AccAddress(tmcrypto.GenPrivKey().PubKey().Address()),
		Secret:          nil,
	}
	res = h(ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, types.AttendeeDoesNotExist, res.Code)

	// 1. invalid secret
	msg = types.MsgClaimKey{
		AttendeeAddress: addr,
		Secret:          nil,
	}
	res = h(ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, types.InvalidCommitmentReveal, res.Code)

	// 3. valid claim
	msg = types.MsgClaimKey{
		AttendeeAddress: addr,
		Secret:          secret,
	}
	res = h(ctx, msg)
	require.True(t, res.IsOK())

	// Check that the attendee was updated accordingly
	a, ok := app.LongyKeeper.GetAttendee(ctx, addr)
	require.True(t, ok)
	require.True(t, a.IsClaimed())
	require.Nil(t, a.CurrentCommitment())

	// 4. Cannot reclaim an attendee in a claimed state
	res = h(ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, types.AttendeeAlreadyClaimed, res.Code)
}
