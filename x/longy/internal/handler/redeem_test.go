package handler_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redeem Handler Tests", func() {

	BeforeEach(func() {
		BeforeTestRun()
		//create public account addresses
		sender = util.IDToAddress(qr1)
		receiver = util.IDToAddress(qr2)
	})

	It("should fail when the sender is not the master account", func() {
		msg := types.NewMsgRedeem(sender, sender)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.InsufficientPrivileges))
	})

	Context("when master account set", func() {
		BeforeEach(func() {
			utils.SetMasterAccount(ctx, keeper, sender)
		})

		It("should fail when the sender is the master account, but attendee doesn't exist", func() {
			msg := types.NewMsgRedeem(sender, receiver)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(types.AttendeeNotFound))
		})

		It("should succeed when the sender is the master account and attendee exist", func() {
			utils.SetMasterAccount(ctx, keeper, sender)
			utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true, false)
			msg := types.NewMsgRedeem(sender, receiver)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
		})

		It("should succeed to set all attendee winnings to claimed", func() {
			utils.SetMasterAccount(ctx, keeper, sender)
			attendee := utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true, false)
			attendee.Winnings = append(attendee.Winnings, types.Win{
				Tier:    types.Tier1,
				Name:    "stuff",
				Claimed: false,
			})

			attendee.Winnings = append(attendee.Winnings, types.Win{
				Tier:    types.Tier2,
				Name:    "stuff",
				Claimed: false,
			})
			keeper.SetAttendee(ctx, &attendee)
			msg := types.NewMsgRedeem(sender, receiver)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))

			attendee, exists := keeper.GetAttendee(ctx, attendee.Address)
			Expect(exists).To(BeTrue())
			Expect(len(attendee.Winnings)).To(Equal(2))
			for _, w := range attendee.Winnings {
				Expect(w.Claimed).To(BeTrue())
			}
		})
	})
})
