package querier_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	q "github.com/eco/longy/x/longy/internal/querier"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	abci "github.com/tendermint/tendermint/abci/types"
)

var _ = Describe("Scan Querier Tests", func() {

	var getScan = func(id string) (scan types.Scan, err sdk.Error) {
		res, err := querier(ctx, []string{q.QueryScans, id}, abci.RequestQuery{})
		if err != nil {
			return
		}
		keeper.Cdc.MustUnmarshalJSON(res, &scan)
		return scan, err
	}

	BeforeEach(func() {
		BeforeTestRun()
	})

	It("should fail when id is empty string", func() {
		_, err := getScan("")
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.ScanNotFound))
	})

	It("should fail when id doesn't have a scan", func() {
		_, err := getScan("013a3ad82dce818adf31f4b6230ad06699ad4f351b633fc9149b30d82f3a6d4463704b62a278f9a1637098")
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.ScanNotFound))
	})

	It("should succeed when id has a scan", func() {
		d1 := []byte{1}
		d2 := []byte{2}
		s1 := util.IDToAddress("1234")
		s2 := util.IDToAddress("asdf")
		var err sdk.Error
		scan, err := types.NewScan(s1, s2, d1, d2, 1, 2)
		Expect(err).To(BeNil())
		keeper.SetScan(ctx, scan)

		qscan, err := getScan(types.Encode(scan.ID))
		Expect(err).To(BeNil())
		Expect(string(qscan.ID)).To(Equal(types.Encode(scan.ID)))
	})
})
