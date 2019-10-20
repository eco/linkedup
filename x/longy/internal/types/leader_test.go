package types_test

import (
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Leader Tests", func() {
	BeforeEach(func() {
	})

	It("should leave both tiers of attendee's empty", func() {
		board := types.NewLeaderBoard(100, nil)
		Expect(len(board.Tier1.Attendees)).To(Equal(0))
		Expect(len(board.Tier2.Attendees)).To(Equal(0))
	})

	It("should only have tier 1 attendee's set if not enough users", func() {
		min := types.LeaderBoardTier1Count - 5
		top := make([]types.Attendee, min)
		board := types.NewLeaderBoard(100, top)
		Expect(len(board.Tier1.Attendees)).To(Equal(min))
		Expect(len(board.Tier2.Attendees)).To(Equal(0))
	})

	It("should only have tier 1 full and tier 2 empty", func() {
		min := types.LeaderBoardTier1Count
		top := make([]types.Attendee, min)
		board := types.NewLeaderBoard(100, top)
		Expect(len(board.Tier1.Attendees)).To(Equal(min))
		Expect(len(board.Tier2.Attendees)).To(Equal(0))
	})

	It("should only have tier 1 full and tier 2 partially filled", func() {
		extra := 5
		min := types.LeaderBoardTier1Count + extra
		top := make([]types.Attendee, min)
		board := types.NewLeaderBoard(100, top)
		Expect(len(board.Tier1.Attendees)).To(Equal(types.LeaderBoardTier1Count))
		Expect(len(board.Tier2.Attendees)).To(Equal(extra))
	})

	It("should only have both tiers filled with attendees when given enough", func() {
		min := types.LeaderBoardCount
		top := make([]types.Attendee, min)
		board := types.NewLeaderBoard(100, top)
		Expect(len(board.Tier1.Attendees)).To(Equal(types.LeaderBoardTier1Count))
		Expect(len(board.Tier2.Attendees)).To(Equal(types.LeaderBoardCount - types.LeaderBoardTier1Count))
	})

	It("should drop attendees if passed more than the leader board max", func() {
		min := types.LeaderBoardCount * 2
		top := make([]types.Attendee, min)
		board := types.NewLeaderBoard(100, top)
		Expect(len(board.Tier1.Attendees)).To(Equal(types.LeaderBoardTier1Count))
		Expect(len(board.Tier2.Attendees)).To(Equal(types.LeaderBoardCount - types.LeaderBoardTier1Count))
	})
})
