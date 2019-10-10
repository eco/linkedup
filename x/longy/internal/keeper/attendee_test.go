package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Attendee Keeper Tests", func() {
	var s1, s2 sdk.AccAddress
	var scan types.Scan
	const (
		qr1 = "1234"
		qr2 = "asdf"
	)
	BeforeEach(func() {
		BeforeTestRun()

		s1 = util.IDToAddress(qr1)
		s2 = util.IDToAddress(qr2)
		var err sdk.Error
		scan, err = types.NewScan(s1, s2, nil, nil)
		Expect(err).To(BeNil())
	})

	Context("For scan events", func() {
		//It("should fail when scan not complete", func() {
		//	err := keeper.AwardScanPoints(ctx, scan)
		//	Expect(err).To(Not(BeNil()))
		//	Expect(err.Code()).To(Equal(types.ScanNotComplete))
		//})

		Context("when a can has been completed by both parties", func() {
			BeforeEach(func() {
			})

			It("should fail when attendee does not exist", func() {
				err := keeper.AwardScanPoints(ctx, scan)
				Expect(err).To(Not(BeNil()))
				Expect(err.Code()).To(Equal(types.AttendeeNotFound))
			})

			It("should succeed when both are attendees", func() {
				a1 := utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false)
				a2 := utils.AddAttendeeToKeeper(ctx, &keeper, qr2, false)

				err := keeper.AwardScanPoints(ctx, scan)
				Expect(err).To(BeNil())
				b1, ok := keeper.GetAttendeeWithID(ctx, qr1)
				Expect(ok).To(BeTrue())
				b2, ok := keeper.GetAttendeeWithID(ctx, qr2)
				Expect(ok).To(BeTrue())
				Expect(a1.Rep + types.ScanAttendeeAwardPoints).To(Equal(b1.Rep))
				Expect(a2.Rep + types.ScanAttendeeAwardPoints).To(Equal(b2.Rep))
			})

			It("should succeed when one is a sponsor", func() {
				a1 := utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false)
				a2 := utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true)

				err := keeper.AwardScanPoints(ctx, scan)
				Expect(err).To(BeNil())
				b1, ok := keeper.GetAttendeeWithID(ctx, qr1)
				Expect(ok).To(BeTrue())
				b2, ok := keeper.GetAttendeeWithID(ctx, qr2)
				Expect(ok).To(BeTrue())
				Expect(a1.Rep + types.ScanSponsorAwardPoints).To(Equal(b1.Rep))
				Expect(a2.Rep + types.ScanSponsorAwardPoints).To(Equal(b2.Rep))
			})
		})
	})
})
