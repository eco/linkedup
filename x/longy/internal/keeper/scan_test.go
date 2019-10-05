package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scan Tests", func() {
	var keeper longy.Keeper
	var scan *types.Scan
	BeforeEach(func() {
		BeforeTestRun()
		keeper = simApp.LongyKeeper
		Expect(keeper).ToNot(BeNil())

		s1 := util.IDToAddress("1234")
		s2 := util.IDToAddress("asdf")
		var err sdk.Error
		scan, err = types.NewScan(s1, s2)
		Expect(err).To(BeNil())
	})

	It("should fail when we try to get a scan that doesn't exist", func() {
		_, err := keeper.GetScanByID(ctx, scan.ID)
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.ScanNotFound))
	})

	It("should succeed to get a scan that exists", func() {
		keeper.SetScan(ctx, scan)

		storeScan, err := keeper.GetScanByID(ctx, scan.ID)
		Expect(err).To(BeNil())
		Expect(storeScan.ID).To(Equal(scan.ID))
		Expect(storeScan.S1.Equals(scan.S1)).To(BeTrue())
		Expect(storeScan.S2.Equals(scan.S2)).To(BeTrue())
		Expect(storeScan.Complete).To(Equal(scan.Complete))
	})
})
