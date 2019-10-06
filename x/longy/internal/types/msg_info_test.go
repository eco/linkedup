package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsgInfo Tests", func() {
	var s1, s2 sdk.AccAddress
	BeforeEach(func() {
		s1 = util.IDToAddress("1234")
		s2 = util.IDToAddress("asdf")
	})

	It("should fail when sender address is not set", func() {
		msg := MsgInfo{
			Sender: sdk.AccAddress{},
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should fail when receiver address is not set", func() {
		msg := MsgInfo{
			Sender:   s1,
			Receiver: sdk.AccAddress{},
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should fail if user is trying to share data with themselves", func() {
		msg := MsgInfo{
			Sender:   s1,
			Receiver: s1,
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(CantShareWithSelf))
	})

	It("should fail when data is not set", func() {
		msg := MsgInfo{
			Sender:   s1,
			Receiver: s2,
			Data:     nil,
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(DataCannotBeEmpty))
	})

	It("should fail when data is above the limit", func() {
		data := make([]byte, MaxDataSize+1)
		msg := MsgInfo{
			Sender:   s1,
			Receiver: s2,
			Data:     data,
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(DataSizeOverLimit))
	})

	It("should succeed when everything is set", func() {
		data := make([]byte, 10)
		msg := MsgInfo{
			Sender:   s1,
			Receiver: s2,
			Data:     data,
		}
		err := msg.ValidateBasic()
		Expect(err).To(BeNil())
	})
})
