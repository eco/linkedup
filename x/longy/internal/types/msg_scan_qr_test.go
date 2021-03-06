package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	. "github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsgScanQr Tests", func() {

	BeforeEach(func() {
	})

	It("should fail when sender address is not set", func() {
		msg := MsgScanQr{
			Sender:    sdk.AccAddress{},
			ScannedQR: "",
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should fail when scanned qr code is not set", func() {
		msg := MsgScanQr{
			Sender:    util.IDToAddress("1234"),
			ScannedQR: "",
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(QRCodeInvalid))
	})

	It("should fail when scanned qr code is invalid", func() {
		msg := MsgScanQr{
			Sender:    util.IDToAddress("1234"),
			ScannedQR: "12s345g678",
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(QRCodeInvalid))
	})

	It("should fail when scanned qr code is negative", func() {
		msg := MsgScanQr{
			Sender:    util.IDToAddress("1234"),
			ScannedQR: "-123456789",
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(QRCodeInvalid))
	})

	It("should fail when scanned qr code is a fraction", func() {
		msg := MsgScanQr{
			Sender:    util.IDToAddress("1234"),
			ScannedQR: "-123.5",
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(QRCodeInvalid))
	})

	It("should successfully validate basic on valid MsgScanQr", func() {
		msg := MsgScanQr{
			Sender:    util.IDToAddress("1234"),
			ScannedQR: "123456789",
		}
		err := msg.ValidateBasic()
		Expect(err).To(BeNil())
	})
})
