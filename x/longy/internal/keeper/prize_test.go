package keeper_test

import (
	"bytes"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prize Keeper Tests", func() {
	var prizes types.GenesisPrizes
	BeforeEach(func() {
		BeforeTestRun()
		prizes = types.GetGenesisPrizes()
		//hard-code this to make sure it doesn't un-expectantly change
		Expect(len(prizes)).To(Equal(5))
	})

	It("should fail when we try to get a prize that doesn't exist", func() {
		p := types.Prize{
			Tier: 10,
		}
		_, err := keeper.GetPrize(ctx, p.GetID())
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.PrizeNotFound))
	})

	It("should succeed to put and get a prize", func() {
		prize := prizes[0]
		keeper.SetPrize(ctx, &prize)

		stored, err := keeper.GetPrize(ctx, prize.GetID())
		Expect(err).To(BeNil())
		comparePrizes(prize, stored)
	})

	It("should succeed to get all prizes that exists", func() {
		for i := range prizes {
			keeper.SetPrize(ctx, &prizes[i])
		}

		for _, prize := range prizes {
			stored, err := keeper.GetPrize(ctx, prize.GetID())
			Expect(err).To(BeNil())
			comparePrizes(prize, stored)
		}
	})
})

func comparePrizes(expected types.Prize, actual types.Prize) {
	Expect(expected.Tier).To(Equal(actual.Tier))
	Expect(expected.PrizeText).To(Equal(actual.PrizeText))
	Expect(expected.Quantity).To(Equal(actual.Quantity))
	Expect(expected.RepNeeded).To(Equal(actual.RepNeeded))
	Expect(bytes.Compare(expected.GetID(), actual.GetID())).To(Equal(0))
}
