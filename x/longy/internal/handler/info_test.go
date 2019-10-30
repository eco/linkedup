package handler_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info Handler Tests", func() {
	var data = []byte{1, 2, 3, 2, 1}

	BeforeEach(func() {
		BeforeTestRun()
		//create public account addresses
		sender = util.IDToAddress(qr1)
		receiver = util.IDToAddress(qr2)
		Expect(receiver).To(Not(BeNil()))
	})

	It("should fail when sender attendee doesn't exist", func() {
		msg := types.NewMsgInfo(sender, receiver, data)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AttendeeNotFound))
	})

	Context("when sender exists", func() {
		BeforeEach(func() {
			utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)
		})

		It("should fail when receiver attendee doesn't exist", func() {
			msg := types.NewMsgInfo(sender, receiver, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(types.AttendeeNotFound))
		})

		Context("when both sender and receiver exist", func() {
			BeforeEach(func() {
				utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true, false)
			})

			It("should fail when there is no scan for them", func() {
				msg := types.NewMsgInfo(sender, receiver, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanNotFound))
			})

			Context("when there is a non-accepted scan", func() {
				BeforeEach(func() {
					createScan(qr1, qr2, sender, receiver, nil, false, false)
					inspectScan(sender, receiver, 0, 0, false)

					//set prizes since they can move tiers on claim for the beta testing
					prizes := types.GetGenesisPrizes()
					for i := range prizes {
						keeper.SetPrize(ctx, &prizes[i])
					}
				})

				It("should not update rep if MsgInfo is sent without data", func() {
					msg := types.NewMsgInfo(sender, receiver, nil)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))
					inspectScan(sender, receiver, 0, 0, false)
				})

				It("should not update rep s1 sends with data", func() {
					msg := types.NewMsgInfo(sender, receiver, data)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))
					inspectScan(sender, receiver, 0, 0, false)
				})

				It("should update scan rep for both s1 and s2 on no data", func() {
					msg := types.NewMsgInfo(receiver, sender, nil)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))
					points := types.ScanAttendeeAwardPoints
					inspectScan(sender, receiver, points, points, true)
				})

				It("should add shared points when s2 accepts", func() {
					msg := types.NewMsgInfo(sender, receiver, data)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))
					inspectScan(sender, receiver, 0, 0, false)

					msg = types.NewMsgInfo(receiver, sender, data)
					result = handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))
					points := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
					inspectScan(sender, receiver, points, points, true)
				})

				It("should add scan and share points for s1 and s2 when s2 accepts and one is a sponsor", func() {
					createScan(qr1, qr2, sender, receiver, data, true, false) //make sponsor
					inspectScan(sender, receiver, 0, 0, false)
					//
					msg := types.NewMsgInfo(receiver, sender, data)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))
					attendeePoints := types.ScanSponsorAwardPoints + types.ShareSponsorAwardPoints
					sponsorPoints := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
					inspectScan(sender, receiver, sponsorPoints, attendeePoints, true)
				})

				Context("when an accepted and full data share scan exists", func() {
					points := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
					BeforeEach(func() {
						//Add the partial scan to the keeper
						createScan(qr1, qr2, sender, receiver, data, false, false)
						inspectScan(sender, receiver, 0, 0, false)
						msg := types.NewMsgInfo(receiver, sender, data)
						result := handler(ctx, msg)
						Expect(result.Code).To(Equal(sdk.CodeOK))
						inspectScan(sender, receiver, points, points, true)
					})

					It("should not allow us to reset data and earn more points", func() {
						//Add the partial scan to the keeper
						msg := types.NewMsgInfo(sender, receiver, nil)
						result := handler(ctx, msg)
						Expect(result.Code).To(Equal(sdk.CodeOK))
						inspectScan(sender, receiver, points, points, true)

						msg = types.NewMsgInfo(receiver, sender, nil)
						result = handler(ctx, msg)
						Expect(result.Code).To(Equal(sdk.CodeOK))
						inspectScan(sender, receiver, points, points, true)

						msg = types.NewMsgInfo(sender, receiver, data)
						result = handler(ctx, msg)
						Expect(result.Code).To(Equal(sdk.CodeOK))
						inspectScan(sender, receiver, points, points, true)

						msg = types.NewMsgInfo(receiver, sender, data)
						result = handler(ctx, msg)
						Expect(result.Code).To(Equal(sdk.CodeOK))
						inspectScan(sender, receiver, points, points, true)
					})
				})

			})
		})
	})
})

//Contains checks to see if a value is in the array
//nolint:unused
func Contains(vals []string, v string) (contains bool) {
	for _, a := range vals {
		if a == v {
			return true
		}
	}
	return false
}
