package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MasterKey Keeper Tests", func() {
	var s1, s2 sdk.AccAddress
	const (
		qr1 = "1234"
		qr2 = "asdf"
	)
	BeforeEach(func() {
		BeforeTestRun()

		s1 = util.IDToAddress(qr1)
		s2 = util.IDToAddress(qr2)
	})

	It("should fail to set an empty AccAddress", func() {
		err := keeper.SetMasterAddress(ctx, sdk.AccAddress{})
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(sdk.CodeInvalidAddress))
	})

	It("should fail to set a valid AccAddress for an account that doesn't exist", func() {
		err := keeper.SetMasterAddress(ctx, s1)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(sdk.CodeUnknownAddress))
	})

	It("should fail to return a master account when its not set", func() {
		is := keeper.IsMasterAccount(ctx, s1)
		Expect(is).To(Equal(false))
	})

	It("should succeed to set a master account", func() {
		acc := keeper.AccountKeeper().NewAccountWithAddress(ctx, s1)
		keeper.AccountKeeper().SetAccount(ctx, acc)
		err := keeper.SetMasterAddress(ctx, s1)
		Expect(err).To(BeNil())
	})

	It("should fail to is master account when passed the wrong AccAddress", func() {
		utils.SetMasterAccount(ctx, keeper, s1)

		acc := keeper.AccountKeeper().NewAccountWithAddress(ctx, s2)
		keeper.AccountKeeper().SetAccount(ctx, acc)
		is := keeper.IsMasterAccount(ctx, s2)
		Expect(is).To(BeFalse())
	})

	It("should succeed to return a master account when set", func() {
		utils.SetMasterAccount(ctx, keeper, s1)

		is := keeper.IsMasterAccount(ctx, s1)
		Expect(is).To(BeTrue())
	})
})
