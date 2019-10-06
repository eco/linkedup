package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info Tests", func() {
	var s1, s2 sdk.AccAddress
	BeforeEach(func() {
		s1 = util.IDToAddress("1234")
		s2 = util.IDToAddress("asdf")
	})

	It("should fail when an address is empty", func() {
		_, err := NewInfo(s1, sdk.AccAddress{}, nil)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(AccountAddressEmpty))
	})

	It("should fail when the addresses are the same", func() {
		_, err := NewInfo(s1, s1, nil)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(AccountsSame))
	})

	It("should fail when data is nil", func() {
		_, err := NewInfo(s1, s2, nil)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(DataCannotBeEmpty))
	})

	It("should fail when data is empty", func() {
		_, err := NewInfo(s1, s2, []byte{})
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(DataCannotBeEmpty))
	})

	It("should succeed when addresses are different and not empty", func() {
		data := []byte{1, 2, 3, 4, 3, 2, 1}
		info, err := NewInfo(s1, s2, data)
		Expect(err).To(BeNil())
		Expect(info).To(Not(BeNil()))
		Expect(bytes.Compare(info.Data, data)).To(Equal(0))
	})
})
