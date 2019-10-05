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

	BeforeEach(func() {
		empty = sdk.AccAddress{}
		s1 = util.IDToAddress("1234")
		s2 = util.IDToAddress("asdf")
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
		Expect(err.Code()).To(Equal(types.ScanAccountsSame))
	})

	It("should return the same id regardless of account params order", func() {
		scan0, err := types.NewScan(s1, s2)
		Expect(err).To(BeNil())
		Expect(len(scan0.ID)).To(Equal(len(s2) * 2))

		scan1, err := types.NewScan(s2, s1)
		Expect(err).To(BeNil())
		Expect(len(scan1.ID)).To(Equal(len(s1) * 2))

		Expect(bytes.Equal(scan0.ID, scan1.ID)).To(BeTrue())
	})
})
