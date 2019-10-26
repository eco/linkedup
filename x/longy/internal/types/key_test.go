package types_test

import (
	"bytes"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Key Tests", func() {
	BeforeEach(func() {
	})

	It("should prefix two arrays", func() {
		prefix := []byte("foo")
		key := []byte("bar")

		prefixedKey := types.PrefixKey(prefix, key)
		Expect(bytes.Compare(prefixedKey, []byte("foo::bar"))).To(Equal(0))
	})

	Context("attendee keys", func() {
		It("should return false when key is nil or empty", func() {

			Expect(types.IsAttendeeKey(nil)).To(BeFalse())
			Expect(types.IsAttendeeKey([]byte{})).To(BeFalse())
		})

		It("should return false when key is not of attendee type", func() {
			key := types.MasterKey()

			Expect(types.IsAttendeeKey(key)).To(BeFalse())

			key = types.PrizeKey([]byte("stuff"))
			Expect(types.IsAttendeeKey(key)).To(BeFalse())
		})

		It("should return true when key is of attendee type", func() {
			s1 := util.IDToAddress("asdf")
			key := types.AttendeeKey(s1)

			Expect(types.IsAttendeeKey(key)).To(BeTrue())
		})
	})

	Context("scan keys", func() {
		It("should return false when key is nil or empty", func() {

			Expect(types.IsScanKey(nil)).To(BeFalse())
			Expect(types.IsScanKey([]byte{})).To(BeFalse())
		})

		It("should return false when key is not of attendee type", func() {
			key := types.MasterKey()

			Expect(types.IsScanKey(key)).To(BeFalse())

			key = types.PrizeKey([]byte("stuff"))
			Expect(types.IsScanKey(key)).To(BeFalse())
		})

		It("should return true when key is of attendee type", func() {
			key := types.ScanKey([]byte("asdf"))

			Expect(types.IsScanKey(key)).To(BeTrue())
		})
	})
})
