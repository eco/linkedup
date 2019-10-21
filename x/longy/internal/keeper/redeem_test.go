package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redeem Keeper Tests", func() {
	var s1 sdk.AccAddress
	const (
		qr1 = "1234"
		qr2 = "asdf"
	)
	BeforeEach(func() {
		BeforeTestRun()

		s1 = util.IDToAddress(qr1)
	})

	Context("when attendees don't exist", func() {

		It("should fail fail to claim prizes", func() {
			err := keeper.RedeemPrizes(ctx, s1)
			Expect(err).To(Not(BeNil()))
			Expect(err.Code()).To(Equal(types.AttendeeNotFound))
		})

		Context("when attendee exists", func() {
			var a types.Attendee
			BeforeEach(func() {
				a = utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)

				prizes := types.GetGenesisPrizes()
				for i := range prizes {
					keeper.SetPrize(ctx, &prizes[i])
				}
			})

			It("should succeed when no prizes for attendee", func() {
				Expect(len(a.Winnings)).To(Equal(0))
				err := keeper.RedeemPrizes(ctx, s1)
				Expect(err).To(BeNil())
			})

			It("should succeed when all winnings are initially unclaimed", func() {
				err := keeper.AddRep(ctx, &a, types.Tier1Rep)
				Expect(err).To(BeNil())
				err = keeper.AddRep(ctx, &a, types.Tier2Rep)
				Expect(err).To(BeNil())
				var exists bool
				a, exists = keeper.GetAttendee(ctx, a.Address)
				Expect(exists).To(BeTrue())
				Expect(len(a.Winnings)).To(Equal(2))
				Expect(a.Winnings[0].Claimed).To(BeFalse())
				Expect(a.Winnings[1].Claimed).To(BeFalse())

				err = keeper.RedeemPrizes(ctx, s1)
				Expect(err).To(BeNil())

				a, exists = keeper.GetAttendee(ctx, a.Address)
				Expect(exists).To(BeTrue())
				Expect(len(a.Winnings)).To(Equal(2))
				Expect(a.Winnings[0].Claimed).To(BeTrue())
				Expect(a.Winnings[1].Claimed).To(BeTrue())
			})

			It("should succeed when there are claimed and unclaimed winnings ", func() {
				winning := &types.Win{
					Tier:    types.Tier1,
					Name:    "Stuff",
					Claimed: true,
				}
				a.Winnings = append(a.Winnings, *winning)
				keeper.SetAttendee(ctx, &a)
				added := a.AddWinning(&types.Win{
					Tier:    types.Tier2,
					Name:    "More Stuff",
					Claimed: false,
				})
				Expect(added).To(BeTrue())
				keeper.SetAttendee(ctx, &a)

				err := keeper.RedeemPrizes(ctx, s1)
				Expect(err).To(BeNil())
				var exists bool
				a, exists = keeper.GetAttendee(ctx, a.Address)
				Expect(exists).To(BeTrue())
				Expect(len(a.Winnings)).To(Equal(2))
				Expect(a.Winnings[0].Claimed).To(BeTrue())
				Expect(a.Winnings[1].Claimed).To(BeTrue())
			})
		})
	})
})
