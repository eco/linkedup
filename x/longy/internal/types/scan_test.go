package types_test

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scan Tests", func() {
	var empty, s1, s2 sdk.AccAddress
	var prefixLen int
	BeforeEach(func() {
		empty = sdk.AccAddress{}
		s1 = util.IDToAddress("1234")
		s2 = util.IDToAddress("asdf")
		prefixLen = len(types.Prefix(types.ScanPrefix))
	})

	It("should fail when any address is empty", func() {
		_, err := types.NewScan(s1, empty)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.AccountAddressEmpty))

		_, err = types.NewScan(empty, s1)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.AccountAddressEmpty))

		_, err = types.NewScan(empty, empty)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.AccountAddressEmpty))
	})

	It("should fail when both addresses are the same", func() {
		_, err := types.NewScan(s2, s2)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.AccountsSame))
	})

	It("should return the same id regardless of account params order", func() {
		scan0, err := types.NewScan(s1, s2)
		Expect(err).To(BeNil())

		scan1, err := types.NewScan(s2, s1)
		Expect(err).To(BeNil())
		Expect(len(scan1.ID)).To(Equal(len(s1)*2 + prefixLen))

		Expect(bytes.Equal(scan0.ID, scan1.ID)).To(BeTrue())
	})

	It("should prefix the id with the scan prefix key", func() {
		scan0, err := types.NewScan(s1, s2)
		keyIncrease := len(s2)*2 + prefixLen
		Expect(err).To(BeNil())
		Expect(len(scan0.ID)).To(Equal(keyIncrease))
		Expect(bytes.Compare(scan0.ID[:1], types.ScanPrefix)).To(Equal(0))
	})
})
