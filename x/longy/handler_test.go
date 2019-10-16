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
	var masterAddr sdk.AccAddress = util.IDToAddress("master")
	BeforeEach(func() {
		/*
		* Here we setup a test attendee & cosmos account that has a
		* id of "1"
		 */
		BeforeTestRun()
		addr = util.IDToAddress("1")

		keeper = simApp.LongyKeeper
		Expect(keeper).ToNot(BeNil())

		handler = longy.NewHandler(keeper)
		Expect(handler).ToNot(BeNil())

		genesis := longy.GenesisState{
			KeyService: types.GenesisKeyService{
				Address: masterAddr,
			},
			Redeem: types.GenesisRedeemKey{
				Address: util.IDToAddress("redeem"),
			},
			Attendees: []types.GenesisAttendee{
				types.GenesisAttendee{
					ID: "1",
				},
			},
		}

		// Init from genesis state
		longy.InitGenesis(ctx, keeper, genesis)
	})

	It("rejects an invalid msg", func() {
		res := handler(ctx, sdk.NewTestMsg())
		Expect(res.IsOK()).Should(BeFalse())
		Expect(strings.Contains(res.Log, "unrecognized longy msg type")).
			Should(BeTrue())
	})

	It("rejets a non-existent attendee", func() {
		addr = util.IDToAddress("2")
		msg := types.MsgKey{
			AttendeeAddress:      addr,
			MasterAddress:        masterAddr,
			NewAttendeePublicKey: nil,
			Commitment:           nil,
		}
		res := handler(ctx, msg)
		Expect(res.IsOK()).Should(BeFalse())
		Expect(res.Code).To(Equal(types.AttendeeNotFound))
	})

	It("checks that the initial public key is nil", func() {
		accountKeeper := keeper.AccountKeeper()
		account := accountKeeper.GetAccount(ctx, addr)
		Expect(account).ToNot(BeNil())

		pubKey := account.GetPubKey()
		Expect(pubKey).To(BeNil())

		attendee, ok := keeper.GetAttendee(ctx, addr)
		Expect(ok).Should(BeTrue())
		Expect(attendee.IsClaimed()).Should(BeFalse())
	})

	It("changes public keys with the key message", func() {
		_, commitment := util.CreateCommitment()
		newPub := tmcrypto.GenPrivKey().PubKey()

		/** setup a key against the account **/
		msg := types.MsgKey{
			AttendeeAddress:      addr,
			MasterAddress:        masterAddr,
			NewAttendeePublicKey: newPub,
			Commitment:           commitment,
		}

		res := handler(ctx, msg)
		Expect(res.IsOK()).Should(BeTrue())

		// Account should have swapped keys
		acc := keeper.AccountKeeper().GetAccount(ctx, addr)
		Expect(acc).ToNot(BeNil())
		Expect(newPub.Address()).To(Equal(acc.GetPubKey().Address()))

		// Attendee must not be in the claimed state
		a, ok := keeper.GetAttendee(ctx, addr)
		Expect(ok).Should(BeTrue())
		Expect(a.IsClaimed()).Should(BeFalse())
		Expect(commitment).To(Equal(a.CurrentCommitment()))
	})

	var _ = Context("with an attendee that has keyed", func() {
		var secret string
		BeforeEach(func() {
			s, commitment := util.CreateCommitment()
			secret = s
			newPub := tmcrypto.GenPrivKey().PubKey()

			/** setup a key against the account **/
			msg := types.MsgKey{
				AttendeeAddress:      addr,
				MasterAddress:        masterAddr,
				NewAttendeePublicKey: newPub,
				Commitment:           commitment,
			}
			res := handler(ctx, msg)
			Expect(res.IsOK()).Should(BeTrue())
		})

		It("cannot key again", func() {
			msg := types.MsgKey{
				AttendeeAddress:      addr,
				MasterAddress:        masterAddr,
				NewAttendeePublicKey: nil,
				Commitment:           nil,
			}
			res := handler(ctx, msg)
			Expect(res.IsOK()).Should(BeFalse())
			Expect(res.Code).Should(Equal(types.AttendeeKeyed))
		})

		It("rejects an invalid commitment", func() {
			msg := types.MsgClaimKey{
				AttendeeAddress: addr,
				Secret:          "",
			}
			res := handler(ctx, msg)
			Expect(res.IsOK()).Should(BeFalse())
			Expect(res.Code).To(Equal(types.InvalidCommitmentReveal))
		})

		It("can claim the account", func() {
			msg := types.MsgClaimKey{
				AttendeeAddress: addr,
				Secret:          secret,
			}
			res := handler(ctx, msg)
			Expect(res.IsOK()).Should(BeTrue())

			// Check that the attendee was updated accordingly
			a, ok := keeper.GetAttendee(ctx, addr)
			Expect(ok).Should(BeTrue())
			Expect(a.IsClaimed()).Should(BeTrue())

			Expect(a.GetRep()).To(Equal(uint(5)))
		})

		It("cannot claim twice", func() {
			msg := types.MsgClaimKey{
				AttendeeAddress: addr,
				Secret:          secret,
			}
			res := handler(ctx, msg)
			Expect(res.IsOK()).Should(BeTrue())

			// try again
			res = handler(ctx, msg)
			Expect(res.IsOK()).Should(BeFalse())
			Expect(res.Code).To(Equal(types.AttendeeClaimed))
		})
	})
})
