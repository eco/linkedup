package keeper_test

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scan Keeper Tests", func() {
	var scan *types.Scan
	d1 := []byte{1}
	d2 := []byte{2}
	BeforeEach(func() {
		BeforeTestRun()

		s1 := util.IDToAddress("1234")
		s2 := util.IDToAddress("asdf")
		var err sdk.Error
		scan, err = types.NewScan(s1, s2, d1, d2, 1, 2)
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
		Expect(bytes.Compare(storeScan.D1, scan.D1)).To(Equal(0))
		Expect(bytes.Compare(storeScan.D2, scan.D2)).To(Equal(0))
		Expect(storeScan.P1).To(Equal(scan.P1))
		Expect(storeScan.P2).To(Equal(scan.P2))
	})
})
