package types_test

import (
	"bytes"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Attendee Tests", func() {

	BeforeEach(func() {
	})

	It("should fail if the adding id is nil or empty", func() {
		attendee := types.NewAttendee("asdf", false)
		added := attendee.AddScanID(nil)
		Expect(added).To(BeFalse())
	})

	It("should succeed if the id is not in the array", func() {
		b := []byte{1, 2, 3}
		attendee := types.NewAttendee("asdf", false)
		added := attendee.AddScanID(b)
		Expect(added).To(BeTrue())
	})

	It("should fail when id already in scan ids", func() {
		b := []byte{1, 2, 3}
		attendee := types.NewAttendee("asdf", false)
		added := attendee.AddScanID(b)
		Expect(added).To(BeTrue())

		added = attendee.AddScanID(b)
		Expect(added).To(BeFalse())
	})

	It("should succeed to store by array", func() {
		s1 := util.IDToAddress("asdf")
		s2 := util.IDToAddress("1234")
		id, err := types.GenScanID(s1, s2)
		Expect(err).To(BeNil())
		attendee := types.NewAttendee("asdf", false)
		added := attendee.AddScanID(id)
		Expect(added).To(BeTrue())

		added = attendee.AddScanID(id)
		Expect(added).To(BeFalse())
	})

	It("should succeed to encode/decode array", func() {
		s1 := util.IDToAddress("asdf")
		s2 := util.IDToAddress("1234")
		id, err := types.GenScanID(s1, s2)
		Expect(err).To(BeNil())

		enc := types.Encode(id)
		back := types.Decode(enc)
		Expect(bytes.Compare(id, back)).To(Equal(0))
	})

	It("should return the correct tier for the attendee", func() {
		attendee := types.NewAttendee("asdf", false)
		Expect(attendee.GetTier()).To(Equal(types.Tier0))

		attendee.Rep = types.Tier1Rep
		Expect(attendee.GetTier()).To(Equal(types.Tier1))

		attendee.Rep = types.Tier2Rep
		Expect(attendee.GetTier()).To(Equal(types.Tier2))

		attendee.Rep = types.Tier3Rep
		Expect(attendee.GetTier()).To(Equal(types.Tier3))

		attendee.Rep = types.Tier4Rep
		Expect(attendee.GetTier()).To(Equal(types.Tier4))

		attendee.Rep = types.Tier5Rep
		Expect(attendee.GetTier()).To(Equal(types.Tier5))
		attendee.Rep = types.Tier5Rep + types.Tier5Rep
		Expect(attendee.GetTier()).To(Equal(types.Tier5))
	})

	It("should refuse to add invalid win to winnings", func() {
		attendee := types.NewAttendee("asdf", false)
		Expect(len(attendee.Winnings)).To(Equal(0))

		w := types.Win{
			Tier:    types.Tier1,
			Name:    "Name and stuff",
			Claimed: true,
		}

		Expect(attendee.AddWinning(&w)).To(BeFalse())
	})

	It("should add valid win to winnings", func() {
		attendee := types.NewAttendee("asdf", false)
		Expect(len(attendee.Winnings)).To(Equal(0))

		w := types.Win{
			Tier:    types.Tier1,
			Name:    "Name and stuff",
			Claimed: false,
		}

		Expect(attendee.AddWinning(&w)).To(BeTrue())
		Expect(len(attendee.Winnings)).To(Equal(1))
	})

	Context("when a winning exists", func() {
		var win types.Win
		var attendee types.Attendee
		BeforeEach(func() {
			attendee = types.NewAttendee("asdf", false)
			win = types.Win{
				Tier:    types.Tier1,
				Name:    "Name and stuff",
				Claimed: false,
			}

			Expect(attendee.AddWinning(&win)).To(BeTrue())
			Expect(len(attendee.Winnings)).To(Equal(1))
		})

		It("should fail to add a second winning of the same tier", func() {
			Expect(attendee.AddWinning(&win)).To(BeFalse())
			Expect(len(attendee.Winnings)).To(Equal(1))
		})

		It("should refuse to claim a winning of a different tier", func() {
			win = types.Win{
				Tier: types.Tier2,
			}

			Expect(attendee.ClaimWinning(win.Tier)).To(BeFalse())
			Expect(attendee.Winnings[0].Claimed).To(BeFalse())
		})

		It("should succeed to claim a winning", func() {
			win = types.Win{
				Tier: types.Tier1,
			}

			Expect(attendee.ClaimWinning(win.Tier)).To(BeTrue())
			Expect(attendee.Winnings[0].Claimed).To(BeTrue())
		})

		It("should fail to claim a winning twice", func() {
			win = types.Win{
				Tier: types.Tier1,
			}

			Expect(attendee.ClaimWinning(win.Tier)).To(BeTrue())
			Expect(attendee.Winnings[0].Claimed).To(BeTrue())

			Expect(attendee.ClaimWinning(win.Tier)).To(BeFalse())
		})

	})
})
