package longy_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
)

var _ = Describe("Longy Handler Tests", func() {
	var handler sdk.Handler
	var keeper longy.Keeper
	var addr sdk.AccAddress
	BeforeEach(func() {
		BeforeTestRun()

		keeper = simApp.LongyKeeper
		Expect(keeper).ToNot(BeNil())

		handler = longy.NewHandler(keeper)
		Expect(handler).ToNot(BeNil())

		// setup the base attendee with a badge of 1
		id := "1"
		addr = util.IDToAddress(id)
		acc := keeper.AccountKeeper().NewAccountWithAddress(ctx, addr)
		attendee := types.NewAttendee(id)

		keeper.AccountKeeper().SetAccount(ctx, acc)
		keeper.SetAttendee(ctx, attendee)
	})

	It("Tests against an invalid msg", func() {
		res := handler(ctx, sdk.NewTestMsg())
		Expect(res.IsOK()).Should(BeFalse())
		Expect(strings.Contains(res.Log, "unrecognized longy msg type")).
			Should(BeTrue())
	})

	It("Swaps public keys with the rekey message", func() {
		_, commitment := util.CreateCommitment()
		newPub := tmcrypto.GenPrivKey().PubKey()

		/** setup a rekey against the account **/
		msg := types.MsgRekey{
			AttendeeAddress:      addr,
			NewAttendeePublicKey: newPub,
			Commitment:           commitment,
		}

		res := handler(ctx, msg)
		Expect(res.IsOK()).Should(BeTrue())

		// Account should swap keys
		acc := keeper.AccountKeeper().GetAccount(ctx, addr)
		Expect(acc).ToNot(BeNil())
		Expect(newPub.Address()).To(Equal(acc.GetPubKey().Address()))

		// Attendee must not be in the claimed state
		a, ok := keeper.GetAttendee(ctx, addr)
		Expect(ok).Should(BeTrue())
		Expect(a.IsClaimed()).Should(BeFalse())
		Expect(commitment).To(Equal(a.CurrentCommitment()))
	})

	It("Can claim an Attendee after Rekeying", func() {
		/** setup a rekey against the account **/
		secret, commitment := util.CreateCommitment()
		newPub := tmcrypto.GenPrivKey().PubKey()

		rekeyMsg := types.MsgRekey{
			AttendeeAddress:      addr,
			NewAttendeePublicKey: newPub,
			Commitment:           commitment,
		}

		res := handler(ctx, rekeyMsg)
		Expect(res.IsOK()).Should(BeTrue())

		/** claim the account **/

		// 1. non-existent attendee
		msg := types.MsgClaimKey{
			AttendeeAddress: sdk.AccAddress(tmcrypto.GenPrivKey().PubKey().Address()),
			Secret:          nil,
		}
		res = handler(ctx, msg)
		Expect(res.IsOK()).Should(BeFalse())
		Expect(res.Code).To(Equal(types.AttendeeDoesNotExist))

		// 1. invalid secret
		msg = types.MsgClaimKey{
			AttendeeAddress: addr,
			Secret:          nil,
		}
		res = handler(ctx, msg)
		Expect(res.IsOK()).Should(BeFalse())
		Expect(res.Code).To(Equal(types.InvalidCommitmentReveal))

		// 3. valid claim
		msg = types.MsgClaimKey{
			AttendeeAddress: addr,
			Secret:          secret,
		}
		res = handler(ctx, msg)
		Expect(res.IsOK()).Should(BeTrue())

		// Check that the attendee was updated accordingly
		a, ok := keeper.GetAttendee(ctx, addr)
		Expect(ok).Should(BeTrue())
		Expect(a.IsClaimed()).Should(BeTrue())
		Expect(a.CurrentCommitment()).Should(BeNil())

		// 4. Cannot reclaim an attendee in a claimed state
		res = handler(ctx, msg)
		Expect(res.IsOK()).Should(BeFalse())
		Expect(res.Code).To(Equal(types.AttendeeAlreadyClaimed))
	})
})
