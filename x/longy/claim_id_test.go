package longy_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/eco/longy/x/longy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Claim Id Tests", func() {
	var handler sdk.Handler
	var keeper longy.Keeper
	BeforeEach(func() {
		BeforeTestRun()
		keeper = simApp.LongyKeeper
		Expect(keeper).ToNot(BeNil())
		handler = longy.NewHandler(keeper)
		Expect(handler).ToNot(BeNil())
	})

	It("should fail unknown message for unsupported type", func() {
		msg := bank.MsgSend{}
		result := handler(ctx, msg)
		Expect(result).ToNot(BeNil())
		Expect(result.Code).To(Equal(types.CodeUnknownRequest))
	})

	It("should fail to claim id when the signer is not the super user", func() {

	})

	It("should fail to claim id when the id is not in the keeper", func() {

	})

	It("should succeed to claim id when the id is in the keeper and the singer is the super user", func() {

	})
})
