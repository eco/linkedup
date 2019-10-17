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
				utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)
				utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true, false)
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
					utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true, true)

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

				Context("when attendee is close to switching tiers", func() {
					var sender, receiver types.Attendee
					var prize1, prize2 types.Prize
					BeforeEach(func() {
						var err sdk.Error
						sender, receiver, err = keeper.GetAttendees(ctx, s1, s2)
						Expect(err).To(BeNil())

						sender.Rep = types.Tier1Rep - 1 - types.ScanAttendeeAwardPoints
						keeper.SetAttendee(ctx, &sender)

						prizes := types.GetGenesisPrizes()
						for i := range prizes {
							keeper.SetPrize(ctx, &prizes[i])
						}
						prize1 = prizes[0]
						prize2 = prizes[1]
					})

					It("should fail to award a prize when attendee does not pass a tier", func() {
						err := keeper.AwardScanPoints(ctx, scan)
						Expect(err).To(BeNil())
						a, ok := keeper.GetAttendee(ctx, s1)
						Expect(ok).To(BeTrue())
						Expect(len(a.Winnings)).To(Equal(0))
					})

					It("should succeed to award a prize when attendee passes a tier", func() {
						sender.Rep = types.Tier1Rep - 1
						keeper.SetAttendee(ctx, &sender)
						receiver.Rep = types.Tier2Rep - 1
						keeper.SetAttendee(ctx, &receiver)

						err := keeper.AwardScanPoints(ctx, scan)
						Expect(err).To(BeNil())
						a, b, err := keeper.GetAttendees(ctx, s1, s2)
						Expect(err).To(BeNil())

						Expect(len(a.Winnings)).To(Equal(1))
						Expect(a.Winnings[0].Tier).To(Equal(types.Tier1))
						Expect(a.Winnings[0].Claimed).To(BeFalse())
						Expect(a.Rep).To(Equal(types.Tier1Rep))

						Expect(len(b.Winnings)).To(Equal(1))
						Expect(b.Winnings[0].Tier).To(Equal(types.Tier2))
						Expect(b.Winnings[0].Claimed).To(BeFalse())
						Expect(b.Rep).To(Equal(types.Tier2Rep))

						p1, err := keeper.GetPrize(ctx, prize1.GetID())
						Expect(err).To(BeNil())
						p2, err := keeper.GetPrize(ctx, prize2.GetID())
						Expect(err).To(BeNil())

						Expect(p1.Quantity).To(Equal(prize1.Quantity - 1))
						Expect(p2.Quantity).To(Equal(prize2.Quantity - 1))

					})

					It("should fail to award a prize when attendee passes a tier, but there are no more "+
						"prizes left", func() {
						prize1.Quantity = 0
						keeper.SetPrize(ctx, &prize1)

						sender.Rep = types.Tier1Rep - 1
						keeper.SetAttendee(ctx, &sender)
						receiver.Rep = types.Tier2Rep - 1
						keeper.SetAttendee(ctx, &receiver)

						err := keeper.AwardScanPoints(ctx, scan)
						Expect(err).To(BeNil())
						a, b, err := keeper.GetAttendees(ctx, s1, s2)
						Expect(err).To(BeNil())

						Expect(len(a.Winnings)).To(Equal(0))
						Expect(a.Rep).To(Equal(types.Tier1Rep))

						Expect(len(b.Winnings)).To(Equal(1))
						Expect(b.Winnings[0].Tier).To(Equal(types.Tier2))
						Expect(b.Winnings[0].Claimed).To(BeFalse())
						Expect(b.Rep).To(Equal(types.Tier2Rep))

						p1, err := keeper.GetPrize(ctx, prize1.GetID())
						Expect(err).To(BeNil())
						p2, err := keeper.GetPrize(ctx, prize2.GetID())
						Expect(err).To(BeNil())

						Expect(p1.Quantity).To(Equal(0))
						Expect(p2.Quantity).To(Equal(prize2.Quantity - 1))
					})
				})
			})
		})
	})

})

//nolint:gocritic
func inspectAttendees(ctx sdk.Context, s1 sdk.AccAddress, s2 sdk.AccAddress, p1 uint, p2 uint) {
	a1, a2, err := keeper.GetAttendees(ctx, s1, s2)
	Expect(err).To(BeNil())
	Expect(a1.Rep).To(Equal(p1))
	Expect(a2.Rep).To(Equal(p2))
}

//nolint:gocritic
func inspectScan(ctx sdk.Context, scanID []byte, p1 uint, p2 uint) {
	scan, err := keeper.GetScanByID(ctx, scanID)
	Expect(err).To(BeNil())
	Expect(scan.P1).To(Equal(p1))
	Expect(scan.P2).To(Equal(p2))
}
