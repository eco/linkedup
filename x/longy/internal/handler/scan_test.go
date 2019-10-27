package handler_test

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scan Handler Tests", func() {

	BeforeEach(func() {
		BeforeTestRun()
		//create public account addresses
		sender = util.IDToAddress(qr1)
		receiver = util.IDToAddress(qr2)
	})

	It("should fail when attendee for qr code does not exist", func() {
		msg := types.NewMsgQrScan(sender, qr2, nil)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AttendeeNotFound))
	})

	It("should fail if the badge hasn't been claimed", func() {
		//add sender to keeper
		utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false, false)

		msg := types.NewMsgQrScan(sender, qr1, nil)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AttendeeClaimed))
	})

	It("should fail if the sender is trying to scan themselves", func() {
		//add sender to keeper
		utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, false)

		msg := types.NewMsgQrScan(sender, qr1, nil)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AccountsSame))
	})

	It("should succeed to create a new scan record without data", func() {
		createScan(qr1, qr2, sender, receiver, nil, false)
		inspectScan(sender, receiver, 0, 0, false)

	})

	It("should succeed to create a new scan record with data", func() {
		data := []byte("asdfasdfa")
		createScan(qr1, qr2, sender, receiver, data, false)
		inspectScan(sender, receiver, 0, 0, false)
	})

	Context("when a partial scan already exists but doesn't have shared info from both parties", func() {
		var data []byte
		BeforeEach(func() {
			//Add the partial scan to the keeper
			createScan(qr1, qr2, sender, receiver, nil, false)
			inspectScan(sender, receiver, 0, 0, false)
			data = []byte("asdfasdfa")

			//set prizes since they can move tiers on claim for the beta testing
			prizes := types.GetGenesisPrizes()
			for i := range prizes {
				keeper.SetPrize(ctx, &prizes[i])
			}
		})

		It("should add info for s1 without increment points", func() {
			//Add the partial scan to the keeper
			msg := types.NewMsgQrScan(sender, qr2, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver, 0, 0, false)
		})

		It("should add scan points for s1 and s2 when s2 accepts", func() {
			msg := types.NewMsgQrScan(receiver, qr1, nil)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			points := types.ScanAttendeeAwardPoints
			inspectScan(sender, receiver, points, points, true)
		})

		It("should add scan points for s1 and s2 when s2 accepts, and share info for s1", func() {
			msg := types.NewMsgQrScan(sender, qr2, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver, 0, 0, false)

			msg = types.NewMsgQrScan(receiver, qr1, nil)
			result = handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			points := types.ScanAttendeeAwardPoints
			inspectScan(sender, receiver, points+types.ShareAttendeeAwardPoints, points, true)
		})

		It("should add scan and share points for s1 and s2 when s2 accepts", func() {
			msg := types.NewMsgQrScan(sender, qr2, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver, 0, 0, false)

			msg = types.NewMsgQrScan(receiver, qr1, data)
			result = handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			points := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
			inspectScan(sender, receiver, points, points, true)
		})

		It("should add scan and share points for s1 and s2 when s2 accepts and one is a sponsor", func() {
			createScan(qr1, qr2, sender, receiver, nil, true) //make sponsor
			msg := types.NewMsgQrScan(sender, qr2, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver, 0, 0, false)

			msg = types.NewMsgQrScan(receiver, qr1, data)
			result = handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			attendeePoints := types.ScanSponsorAwardPoints + types.ShareSponsorAwardPoints
			sponsorPoints := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
			inspectScan(sender, receiver, sponsorPoints, attendeePoints, true)
		})

		It("should allow s1 to add scan points after acceptance by s2", func() {
			sumPoints := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
			//Add the partial scan to the keeper
			msg := types.NewMsgQrScan(receiver, qr1, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver, types.ScanAttendeeAwardPoints, sumPoints, true)

			msg = types.NewMsgQrScan(sender, qr2, data)
			result = handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver, sumPoints, sumPoints, true)
		})

		Context("when an accepted and full data share scan exists", func() {
			points := types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints
			BeforeEach(func() {
				//Add the partial scan to the keeper
				createScan(qr1, qr2, sender, receiver, data, false)
				inspectScan(sender, receiver, 0, 0, false)
				msg := types.NewMsgQrScan(receiver, qr1, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))
				inspectScan(sender, receiver, points, points, true)
			})

			It("should not allow us to reset data and earn more points", func() {
				//attendeePoints := types.ScanSponsorAwardPoints + types.ShareSponsorAwardPoints
				////Add the partial scan to the keeper
				msg := types.NewMsgQrScan(sender, qr2, nil)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))
				inspectScan(sender, receiver, points, points, true)

				msg = types.NewMsgQrScan(receiver, qr1, nil)
				result = handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))
				inspectScan(sender, receiver, points, points, true)

				msg = types.NewMsgQrScan(sender, qr2, data)
				result = handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))
				inspectScan(sender, receiver, points, points, true)

				msg = types.NewMsgQrScan(receiver, qr1, data)
				result = handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))
				inspectScan(sender, receiver, points, points, true)
			})
		})
	})
})

func inspectScan(s1 sdk.AccAddress, s2 sdk.AccAddress, p1 uint, p2 uint, accepted bool) {
	id, err := types.GenScanID(s2, s1) //invert for fun, order shouldn't matter
	Expect(err).To(BeNil())
	scan, err := keeper.GetScanByID(ctx, id)
	Expect(err).To(BeNil())
	Expect(scan).To(Not(BeNil()))
	Expect(scan.S1.Equals(s1)).To(BeTrue())
	Expect(scan.S2.Equals(s2)).To(BeTrue())
	Expect(bytes.Compare(scan.ID, id)).To(Equal(0))

	//get attendees
	a, ok := keeper.GetAttendee(ctx, s1)
	Expect(ok).To(BeTrue())
	b, ok := keeper.GetAttendee(ctx, s2)
	Expect(ok).To(BeTrue())

	//check the rep on both the scan and attendees
	Expect(a.Rep).To(Equal(p1))
	Expect(b.Rep).To(Equal(p2))
	Expect(scan.P1).To(Equal(p1))
	Expect(scan.P2).To(Equal(p2))
	//check accepted
	Expect(scan.Accepted).To(Equal(accepted))

	//Check share ids
	Expect(len(a.ScanIDs)).To(Equal(1))
	Expect(len(b.ScanIDs)).To(Equal(1))
	//check actual id
	Expect(a.ScanIDs[0]).To(Equal(b.ScanIDs[0]))
	id, err = types.GenScanID(a.Address, b.Address)
	Expect(err).To(BeNil())
	Expect(bytes.Compare(id, types.Decode(a.ScanIDs[0]))).To(Equal(0))
}

//nolint:unparam
func createScan(qr1 string, qr2 string,
	s1 sdk.AccAddress, s2 sdk.AccAddress, data []byte, sponsor bool) {
	//add sender to keeper
	utils.AddAttendeeToKeeper(ctx, &keeper, qr1, true, sponsor)
	//add scanned to keeper
	utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true, false)

	msg := types.NewMsgQrScan(s1, qr2, data)
	result := handler(ctx, msg)
	Expect(result.Code).To(Equal(sdk.CodeOK))

	//check that D1 is being set
	id, err := types.GenScanID(sender, receiver)
	Expect(err).To(BeNil())
	scan, err := keeper.GetScanByID(ctx, id)
	Expect(err).To(BeNil())
	if data == nil {
		Expect(len(scan.D1)).To(Equal(0))
	} else {
		Expect(bytes.Compare(scan.D1, data)).To(Equal(0))
	}
}
