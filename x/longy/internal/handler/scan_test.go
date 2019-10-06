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
	var s1, s2 sdk.AccAddress
	const (
		qr1 = "1234"
		qr2 = "asdf"
	)
	BeforeEach(func() {
		BeforeTestRun()
		//create public account addresses
		s1 = util.IDToAddress(qr1)
		s2 = util.IDToAddress(qr2)
	})

	It("should fail when attendee for qr code does not exist", func() {
		msg := types.NewMsgQrScan(s1, qr2)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AttendeeNotFound))
	})

	It("should fail if the sender is trying to scan themselves", func() {
		//add sender to keeper
		utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false)

		msg := types.NewMsgQrScan(s1, qr1)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AccountsSame))
	})

	It("should succeed to create a new scan record on first scan", func() {
		createScan(qr1, qr2, s1, s2, false)
	})

	Context("when a partial scan already exists but hasn't been completed by other party", func() {
		BeforeEach(func() {
			//Add the partial scan to the keeper
			createScan(qr1, qr2, s1, s2, false)
		})

		It("should fail if the sender is trying to rescan the same person", func() {
			msg := types.NewMsgQrScan(s1, qr2)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(types.ScanQRAlreadyOccurred))
		})

		It("should succeed if the scanned attendee post their tx in response", func() {
			attendee, ok := keeper.GetAttendee(ctx, s1)
			Expect(ok).To(BeTrue())
			rep := attendee.Rep

			msg := types.NewMsgQrScan(s2, qr1)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))

			inspectScan(s1, s2, true)

			//check rep updated on completion
			attendee, ok = keeper.GetAttendee(ctx, s1)
			Expect(ok).To(BeTrue())
			Expect(attendee.Rep).To(Equal(rep + types.ScanAttendeeAwardPoints))
		})

		It("should succeed if the scanned attendee post their tx in response and one is a sponsor", func() {
			utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true)
			attendee, ok := keeper.GetAttendee(ctx, s1)
			Expect(ok).To(BeTrue())
			rep := attendee.Rep

			msg := types.NewMsgQrScan(s2, qr1)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))

			inspectScan(s1, s2, true)

			//check rep updated on completion
			attendee, ok = keeper.GetAttendee(ctx, s1)
			Expect(ok).To(BeTrue())
			Expect(attendee.Rep).To(Equal(rep + types.ScanSponsorAwardPoints))
		})

		Context("when a scan has been completed by both parties", func() {
			BeforeEach(func() {
				//complete the scan
				msg := types.NewMsgQrScan(s2, qr1)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))

				inspectScan(s1, s2, true)
			})

			It("should fail if called again by the scanner", func() {
				msg := types.NewMsgQrScan(s1, qr2)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanQRAlreadyOccurred))
			})

			It("should fail if called again by the person scanned", func() {
				msg := types.NewMsgQrScan(s2, qr1)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanQRAlreadyOccurred))
			})
		})

	})
})

func inspectScan(s1 sdk.AccAddress, s2 sdk.AccAddress, completed bool) {
	id, err := types.GenScanID(s2, s1) //invert for fun, order shouldn't matter
	Expect(err).To(BeNil())
	scan, err := keeper.GetScanByID(ctx, id)
	Expect(err).To(BeNil())
	Expect(scan).To(Not(BeNil()))
	Expect(scan.S1.Equals(s1)).To(BeTrue())
	Expect(scan.S2.Equals(s2)).To(BeTrue())
	Expect(bytes.Compare(scan.ID, id)).To(Equal(0))
	if completed {
		Expect(scan.Complete).To(BeTrue())
	} else {
		Expect(scan.Complete).To(BeFalse())
	}
}

func createScan(qr1 string, qr2 string,
	s1 sdk.AccAddress, s2 sdk.AccAddress, sponsor bool) {
	//add sender to keeper
	utils.AddAttendeeToKeeper(ctx, &keeper, qr1, sponsor)
	//add scanned to keeper
	utils.AddAttendeeToKeeper(ctx, &keeper, qr2, false)

	msg := types.NewMsgQrScan(s1, qr2)
	result := handler(ctx, msg)
	Expect(result.Code).To(Equal(sdk.CodeOK))
	inspectScan(s1, s2, false)
}
