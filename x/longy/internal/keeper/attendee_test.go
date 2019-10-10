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
	var scan *types.Scan
	const (
		qr1 = "1234"
		qr2 = "asdf"
	)
	BeforeEach(func() {
		BeforeTestRun()

		s1 = util.IDToAddress(qr1)
		s2 = util.IDToAddress(qr2)
		var err sdk.Error
		scan, err = types.NewScan(s1, s2, nil, nil, 0, 0)
		Expect(err).To(BeNil())
	})

	Context("For scan events", func() {
		It("should fail when scan not complete", func() {
			err := keeper.AwardScanPoints(ctx, scan)
			Expect(err).To(Not(BeNil()))
			Expect(err.Code()).To(Equal(types.ScanNotAccepted))
		})

		Context("when a can has been completed by both parties", func() {
			BeforeEach(func() {
				scan.Accepted = true
				keeper.SetScan(ctx, scan)
			})

			It("should fail when attendee does not exist", func() {
				err := keeper.AwardScanPoints(ctx, scan)
				Expect(err).To(Not(BeNil()))
				Expect(err.Code()).To(Equal(types.AttendeeNotFound))

				inspectScan(ctx, scan.ID, 0, 0)
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

				inspectScan(ctx, scan.ID, b1.Rep, b2.Rep)
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
				Expect(a2.Rep + types.ScanAttendeeAwardPoints).To(Equal(b2.Rep))

				inspectScan(ctx, scan.ID, b1.Rep, b2.Rep)
			})
		})
	})
})

//nolint:gocritic
func inspectScan(ctx sdk.Context, scanID []byte, p1 uint, p2 uint) {
	scan, err := keeper.GetScanByID(ctx, scanID)
	Expect(err).To(BeNil())
	Expect(scan.P1).To(Equal(p1))
	Expect(scan.P2).To(Equal(p2))
}
