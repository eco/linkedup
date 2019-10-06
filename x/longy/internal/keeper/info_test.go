package keeper_test

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info Keeper Tests", func() {
	var info types.Info
	BeforeEach(func() {
		BeforeTestRun()

		s1 := util.IDToAddress("1234")
		s2 := util.IDToAddress("asdf")
		var err sdk.Error
		info, err = types.NewInfo(s1, s2, []byte{1, 2, 3})
		Expect(err).To(BeNil())
	})

	It("should fail when we try to get an info that doesn't exist", func() {
		_, err := keeper.GetInfoByID(ctx, info.ID)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.InfoNotFound))
	})

	It("should succeed to get a scan that exists", func() {
		keeper.SetInfo(ctx, &info)

		storeInfo, err := keeper.GetInfoByID(ctx, info.ID)
		Expect(err).To(BeNil())
		Expect(storeInfo.ID).To(Equal(info.ID))
		Expect(bytes.Compare(storeInfo.Data, info.Data)).To(Equal(0))
	})
})
