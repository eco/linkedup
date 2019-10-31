package querier_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	q "github.com/eco/longy/x/longy/internal/querier"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	abci "github.com/tendermint/tendermint/abci/types"
)

var _ = Describe("Winnings Querier Tests", func() {

	var getWinnings = func(addr string) ([]types.Win, sdk.Error) {
		res, err := querier(ctx, []string{q.WinningsKey, addr}, abci.RequestQuery{})
		if err != nil {
			return nil, err
		}
		var winnings []types.Win
		keeper.Cdc.MustUnmarshalJSON(res, &winnings)
		return winnings, err
	}

	BeforeEach(func() {
		BeforeTestRun()
	})

	It("should fail when address is empty string", func() {
		winnings, err := getWinnings(sdk.AccAddress{}.String())
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.AttendeeNotFound))
		Expect(len(winnings)).To(Equal(0))
	})

	It("should fail when address is not bech32", func() {
		winnings, err := getWinnings("asdfBbbb")
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(sdk.CodeInvalidAddress))
		Expect(len(winnings)).To(Equal(0))
	})

	It("should fail when attendee doesn't exist for address", func() {
		sender = util.IDToAddress(qr1)
		winnings, err := getWinnings(sender.String())
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.AttendeeNotFound))
		Expect(len(winnings)).To(Equal(0))
	})

	It("should return an empty array for an attendee that has no winnings", func() {
		attendee := utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)
		winnings, err := getWinnings(attendee.Address.String())
		Expect(err).To(BeNil())
		Expect(len(winnings)).To(Equal(0))
	})

	Context("when there are prizes", func() {
		BeforeEach(func() {
			prizes := types.GetGenesisPrizes()
			for i := range prizes {
				keeper.SetPrize(ctx, &prizes[i])
			}
		})

		It("should return all the unclaimed winnings for an attendee", func() {
			a := utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)
			err := keeper.AddRep(ctx, &a, types.Tier3Rep)
			Expect(err).To(BeNil())
			var exists bool
			a, exists = keeper.GetAttendee(ctx, a.Address)
			Expect(exists).To(BeTrue())
			Expect(len(a.Winnings)).To(Equal(3))
			Expect(a.Winnings[0].Claimed).To(BeFalse())
			Expect(a.Winnings[1].Claimed).To(BeFalse())
			Expect(a.Winnings[2].Claimed).To(BeFalse())

			winnings, err := getWinnings(a.Address.String())
			Expect(err).To(BeNil())
			Expect(len(winnings)).To(Equal(3))
			compare(a.Winnings[0], winnings[0])
			compare(a.Winnings[1], winnings[1])
			compare(a.Winnings[2], winnings[2])
		})

		It("should return claimed and unclaimed winnings for an attendee", func() {
			a := utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)
			err := keeper.AddRep(ctx, &a, types.Tier3Rep)
			Expect(err).To(BeNil())
			var exists bool
			a, exists = keeper.GetAttendee(ctx, a.Address)
			Expect(exists).To(BeTrue())
			Expect(len(a.Winnings)).To(Equal(3))
			Expect(a.Winnings[0].Claimed).To(BeFalse())
			Expect(a.Winnings[1].Claimed).To(BeFalse())
			Expect(a.Winnings[2].Claimed).To(BeFalse())

			a.Winnings[0].Claimed = true
			keeper.SetAttendee(ctx, &a)

			winnings, err := getWinnings(a.Address.String())
			Expect(err).To(BeNil())
			Expect(len(winnings)).To(Equal(3))
			compare(a.Winnings[0], winnings[0])
			compare(a.Winnings[1], winnings[1])
			compare(a.Winnings[2], winnings[2])
		})

	})
})

func compare(expected types.Win, actual types.Win) {
	Expect(expected.Claimed).To(Equal(actual.Claimed))
	Expect(expected.Tier).To(Equal(actual.Tier))
	Expect(expected.Name).To(Equal(actual.Name))
}
