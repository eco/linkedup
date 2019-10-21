package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsgRedeem Tests", func() {

	BeforeEach(func() {
	})

	It("should fail when sender address is not set", func() {
		msg := MsgRedeem{
			Sender: sdk.AccAddress{},
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should fail when attendee address is not set", func() {
		msg := MsgRedeem{
			Sender:   util.IDToAddress("1234"),
			Attendee: sdk.AccAddress{},
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should successfully validate basic on valid MsgRedeem", func() {
		msg := MsgRedeem{
			Sender:   util.IDToAddress("1234"),
			Attendee: util.IDToAddress("asdf"),
		}
		err := msg.ValidateBasic()
		Expect(err).To(BeNil())
	})
})
