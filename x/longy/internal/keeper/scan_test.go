package keeper_test

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
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
		compare(scan, storeScan)
	})

	It("should return no scans when there aren't any", func() {
		scans := keeper.GetAllScans(ctx)
		Expect(len(scans)).To(Equal(0))
	})

	Context("multiple scans exist", func() {
		var scan1, scan2 *types.Scan
		BeforeEach(func() {
			var err error
			scan1, err = types.NewScan(util.IDToAddress("2222"), util.IDToAddress("33333"),
				[]byte{3}, []byte{4}, 5, 6)
			Expect(err).To(BeNil())

			scan2, err = types.NewScan(util.IDToAddress("4444"), util.IDToAddress("55555"),
				[]byte{4}, []byte{5}, 6, 7)
			Expect(err).To(BeNil())

			scan.UnixTimeSec = 1
			scan1.UnixTimeSec = 2
			scan2.UnixTimeSec = 3
			keeper.SetScan(ctx, scan)
			keeper.SetScan(ctx, scan1)
			keeper.SetScan(ctx, scan2)
		})

		It("should return all scans", func() {
			scans := keeper.GetAllScans(ctx)
			Expect(len(scans)).To(Equal(3))
			sort.Slice(scans, func(i, j int) bool { return scans[i].UnixTimeSec < scans[j].UnixTimeSec })
			compare(scan, &scans[0])
			compare(scan1, &scans[1])
			compare(scan2, &scans[2])
		})

		It("should return all scans while ignoring anything else in the keeper", func() {
			keeper.Set(ctx, []byte("somekey"), []byte("somevalue"))
			utils.SetMasterAccount(ctx, keeper, util.IDToAddress("111111"))
			utils.AddAttendeeToKeeper(ctx, &keeper, "111111", true, false)

			scans := keeper.GetAllScans(ctx)
			Expect(len(scans)).To(Equal(3))
		})
	})

})

func compare(expected *types.Scan, actual *types.Scan) {
	Expect(expected.ID).To(Equal(actual.ID))
	Expect(expected.S1.Equals(actual.S1)).To(BeTrue())
	Expect(expected.S2.Equals(actual.S2)).To(BeTrue())
	Expect(expected.P1).To(Equal(actual.P1))
	Expect(expected.P2).To(Equal(actual.P2))
	Expect(bytes.Compare(expected.D1, actual.D1)).To(Equal(0))
	Expect(bytes.Compare(expected.D2, actual.D2)).To(Equal(0))
}
