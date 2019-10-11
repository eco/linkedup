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

	Context("when attendees don't exist", func() {
		It("should fail to award scan points", func() {
			err := keeper.AwardScanPoints(ctx, scan)
			Expect(err).To(Not(BeNil()))
			Expect(err.Code()).To(Equal(types.AttendeeNotFound))
		})

		It("should fail to award share info points", func() {
			err := keeper.AwardShareInfoPoints(ctx, scan, s1, s2)
			Expect(err).To(Not(BeNil()))
			Expect(err.Code()).To(Equal(types.AttendeeNotFound))
		})

		It("should fail add scan ids to participants", func() {
			err := keeper.AddSharedID(ctx, s1, s2, scan.ID)
			Expect(err).To(Not(BeNil()))
			Expect(err.Code()).To(Equal(types.AttendeeNotFound))
		})

		Context("when attendee's exist but a scan has not been accepted", func() {
			BeforeEach(func() {
				BeforeTestRun()
				utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false)
				utils.AddAttendeeToKeeper(ctx, &keeper, qr2, false)
			})

			It("should fail to award scan points", func() {
				err := keeper.AwardScanPoints(ctx, scan)
				Expect(err).To(Not(BeNil()))
				Expect(err.Code()).To(Equal(types.ScanNotAccepted))
				inspectAttendees(ctx, s1, s2, 0, 0)
			})

			It("should fail to award share info points", func() {
				err := keeper.AwardShareInfoPoints(ctx, scan, s1, s2)
				Expect(err).To(Not(BeNil()))
				Expect(err.Code()).To(Equal(types.ScanNotAccepted))
				inspectAttendees(ctx, s1, s2, 0, 0)
			})

			It("should succeed add scan ids to participants", func() {
				err := keeper.AddSharedID(ctx, s1, s2, scan.ID)
				Expect(err).To(BeNil())
				sender, receiver, err := keeper.GetAttendees(ctx, s1, s2)
				Expect(err).To(BeNil())
				Expect(len(sender.ScanIDs)).To(Equal(1))
				Expect(len(receiver.ScanIDs)).To(Equal(1))
				Expect(receiver.ScanIDs[0]).To(Equal(receiver.ScanIDs[0]))
				Expect(receiver.ScanIDs[0]).To(Equal(types.Encode(scan.ID)))
				inspectAttendees(ctx, s1, s2, 0, 0)
			})

			Context("when a scan has been completed by both parties", func() {
				BeforeEach(func() {
					scan.Accepted = true
					keeper.SetScan(ctx, scan)
				})

				It("should succeed to award scan points to attendees", func() {
					err := keeper.AwardScanPoints(ctx, scan)
					Expect(err).To(BeNil())
					point := types.ScanAttendeeAwardPoints
					inspectScan(ctx, scan.ID, point, point)
					inspectAttendees(ctx, s1, s2, point, point)
				})

				It("should succeed to award share info points to attendees", func() {
					err := keeper.AwardShareInfoPoints(ctx, scan, s1, s2)
					Expect(err).To(BeNil())
					point := types.ShareAttendeeAwardPoints
					inspectScan(ctx, scan.ID, point, 0)
					inspectAttendees(ctx, s1, s2, point, 0)

					err = keeper.AwardShareInfoPoints(ctx, scan, s2, s1)
					Expect(err).To(BeNil())
					point = types.ShareAttendeeAwardPoints
					inspectScan(ctx, scan.ID, point, point)
					inspectAttendees(ctx, s1, s2, point, point)
				})

				It("should succeed to award share info points to attendee and sponsor", func() {
					utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true)

					err := keeper.AwardShareInfoPoints(ctx, scan, s1, s2)
					Expect(err).To(BeNil())
					pointSpon := types.ShareSponsorAwardPoints
					inspectScan(ctx, scan.ID, pointSpon, 0)
					inspectAttendees(ctx, s1, s2, pointSpon, 0)

					err = keeper.AwardShareInfoPoints(ctx, scan, s2, s1)
					Expect(err).To(BeNil())
					pointAtt := types.ShareAttendeeAwardPoints
					inspectScan(ctx, scan.ID, pointSpon, pointAtt)
					inspectAttendees(ctx, s1, s2, pointSpon, pointAtt)
				})
			})
		})
	})

})

//nolint:gocritic
func inspectAttendees(ctx sdk.Context, s1 sdk.AccAddress, s2 sdk.AccAddress, p1 uint, p2 uint) {
	a1, a2, err := keeper.GetAttendees(ctx, s1, s2)
	Expect(err).To(BeNil())
	a1.Rep = p1
	a2.Rep = p2
}

//nolint:gocritic
func inspectScan(ctx sdk.Context, scanID []byte, p1 uint, p2 uint) {
	scan, err := keeper.GetScanByID(ctx, scanID)
	Expect(err).To(BeNil())
	Expect(scan.P1).To(Equal(p1))
	Expect(scan.P2).To(Equal(p2))
}
