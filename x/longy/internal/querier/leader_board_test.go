package querier_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy"
	querier2 "github.com/eco/longy/x/longy/internal/querier"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/rand"
)

var _ = FDescribe("Leader Board Querier Tests", func() {

	var getLead = func() types.LeaderBoard {
		res, err := querier(ctx, []string{querier2.LeaderKey}, abci.RequestQuery{})
		Expect(err).To(BeNil())
		var leaderBoard types.LeaderBoard
		keeper.Cdc.MustUnmarshalJSON(res, &leaderBoard)
		return leaderBoard
	}

	BeforeEach(func() {
		BeforeTestRun()
	})

	It("should return empty leader board if no attendees", func() {
		board := getLead()
		Expect(board.TotalCount).To(Equal(0))
		Expect(len(board.Tier1.Attendees)).To(Equal(0))
		Expect(len(board.Tier2.Attendees)).To(Equal(0))
		Expect(board.Tier1.PrizeAmount).To(Equal(types.LeaderBoardTier1Prize))
		Expect(board.Tier2.PrizeAmount).To(Equal(types.LeaderBoardTier2Prize))
	})

	It("should return no more leader board members than attendees in event", func() {
		count := types.LeaderBoardTier1Count - 1
		AddAttendeesToKeeper(ctx, &keeper, count, true)

		board := getLead()
		Expect(board.TotalCount).To(Equal(count))
		Expect(len(board.Tier1.Attendees)).To(Equal(count))
		Expect(len(board.Tier2.Attendees)).To(Equal(0))
		Expect(board.Tier1.PrizeAmount).To(Equal(types.LeaderBoardTier1Prize))
		Expect(board.Tier2.PrizeAmount).To(Equal(types.LeaderBoardTier2Prize))
	})

	It("should return not return attendees that have no rep", func() {
		countWithRep := types.LeaderBoardTier1Count
		countWithOutRep := types.LeaderBoardCount - types.LeaderBoardTier1Count
		AddAttendeesToKeeper(ctx, &keeper, countWithRep, true)
		AddAttendeesToKeeper(ctx, &keeper, countWithOutRep, false)

		attendees := keeper.GetAllAttendees(ctx)
		Expect(len(attendees)).To(Equal(countWithRep + countWithOutRep))

		board := getLead()
		Expect(board.TotalCount).To(Equal(countWithRep + countWithOutRep))
		Expect(len(board.Tier1.Attendees)).To(Equal(countWithRep))
		Expect(len(board.Tier2.Attendees)).To(Equal(0))
		Expect(board.Tier1.PrizeAmount).To(Equal(types.LeaderBoardTier1Prize))
		Expect(board.Tier2.PrizeAmount).To(Equal(types.LeaderBoardTier2Prize))
	})

	It("should return the full leader board", func() {
		count := types.LeaderBoardCount
		AddAttendeesToKeeper(ctx, &keeper, count, true)

		board := getLead()
		Expect(board.TotalCount).To(Equal(count))
		Expect(len(board.Tier1.Attendees)).To(Equal(types.LeaderBoardTier1Count))
		Expect(len(board.Tier2.Attendees)).To(Equal(types.LeaderBoardCount - types.LeaderBoardTier1Count))
		Expect(board.Tier1.PrizeAmount).To(Equal(types.LeaderBoardTier1Prize))
		Expect(board.Tier2.PrizeAmount).To(Equal(types.LeaderBoardTier2Prize))
	})

	It("should not have more attendees in all tiers than the max", func() {
		count := types.LeaderBoardCount * 2
		AddAttendeesToKeeper(ctx, &keeper, count, true)

		board := getLead()
		Expect(board.TotalCount).To(Equal(count))
		Expect(len(board.Tier1.Attendees)).To(Equal(types.LeaderBoardTier1Count))
		Expect(len(board.Tier2.Attendees)).To(Equal(types.LeaderBoardCount - types.LeaderBoardTier1Count))
		Expect(board.Tier1.PrizeAmount).To(Equal(types.LeaderBoardTier1Prize))
		Expect(board.Tier2.PrizeAmount).To(Equal(types.LeaderBoardTier2Prize))
	})

	It("should return order the tiers and attendees in descending order by Rep", func() {
		count := types.LeaderBoardCount * 2
		AddAttendeesToKeeper(ctx, &keeper, count, true)

		board := getLead()
		Expect(board.TotalCount).To(Equal(count))
		u := board.Tier1.Attendees[0].Rep
		for _, a := range board.Tier1.Attendees {
			Expect(u >= a.Rep).To(BeTrue())
			u = a.Rep
		}
		u = board.Tier2.Attendees[0].Rep
		for _, a := range board.Tier2.Attendees {
			Expect(u >= a.Rep).To(BeTrue())
			u = a.Rep
		}
	})
})

//AddAttendeesToKeeper creates the given number of attendees, sets them to claimed and a random Rep between [0,100]
//nolint:gocritic
func AddAttendeesToKeeper(ctx sdk.Context, keeper *longy.Keeper, count int, rep bool) {
	for i := 0; i < count; i++ {
		a := utils.AddAttendeeToKeeper(ctx, keeper, fmt.Sprintf("%d", rand.Int()),
			true, false)
		if rep {
			a.Rep = uint(rand.Intn(100)) + 10
		}
		keeper.SetAttendee(ctx, &a)
	}
}
