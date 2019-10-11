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
		attendee := types.NewAttendee("asdf")
		added := attendee.AddScanID(nil)
		Expect(added).To(BeFalse())
	})

	It("should succeed if the id is not in the array", func() {
		b := []byte{1, 2, 3}
		attendee := types.NewAttendee("asdf")
		added := attendee.AddScanID(b)
		Expect(added).To(BeTrue())
	})

	It("should fail when id already in scan ids", func() {
		b := []byte{1, 2, 3}
		attendee := types.NewAttendee("asdf")
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
		attendee := types.NewAttendee("asdf")
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
		attendee := types.NewAttendee("asdf")
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
})
