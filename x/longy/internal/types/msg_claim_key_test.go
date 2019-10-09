package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("MsgClaimKey Tests", func() {
	var s1 sdk.AccAddress
	var secret, rsa string
	BeforeEach(func() {
		s1 = util.IDToAddress("1234")
		secret = "soccer"
		rsa = "----- Begin Public Key ------"
	})

	It("should fail when attendee address is not set", func() {
		msg := types.MsgClaimKey{}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should fail when secret is not set", func() {
		msg := types.MsgClaimKey{
			AttendeeAddress: s1,
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(types.EmptySecret))
	})

	It("should fail when rsa key is not set", func() {
		msg := types.MsgClaimKey{
			AttendeeAddress: s1,
			Secret:          secret,
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(types.EmptyRsaKey))
	})

	It("should fail when encrypted info is not set", func() {
		msg := types.MsgClaimKey{
			AttendeeAddress: s1,
			Secret:          secret,
			RsaPublicKey:    rsa,
		}
		err := msg.ValidateBasic()
		Expect(err.Error()).To(Not(BeNil()))
		Expect(err.Result().Code).To(Equal(types.EmptyEncryptedInfo))
	})

	It("should succeed when everything is set", func() {
		msg := types.MsgClaimKey{
			AttendeeAddress: s1,
			Secret:          secret,
			RsaPublicKey:    rsa,
			EncryptedInfo:   []byte{1, 2, 3, 4, 5},
		}
		err := msg.ValidateBasic()
		Expect(err).To(BeNil())
	})
})
