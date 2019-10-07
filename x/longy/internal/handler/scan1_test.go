package handler_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = PDescribe("Scan Handler Tests", func() {
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

	Context("when a partial scan already exists but hasn't been completed by other party", func() {
		BeforeEach(func() {
			//Add the partial scan to the keeper
			createScan(qr1, qr2, s1, s2, nil, false)
		})

		It("should fail if the sender is trying to rescan the same person", func() {
			msg := types.NewMsgQrScan(s1, qr2, nil)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(types.ScanQRAlreadyOccurred))
		})

		It("should succeed if the scanned attendee post their tx in response", func() {
			attendee, ok := keeper.GetAttendee(ctx, s1)
			Expect(ok).To(BeTrue())
			rep := attendee.Rep

			msg := types.NewMsgQrScan(s2, qr1, nil)
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

			msg := types.NewMsgQrScan(s2, qr1, nil)
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
				msg := types.NewMsgQrScan(s2, qr1, nil)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))

				inspectScan(s1, s2, true)
			})

			It("should fail if called again by the scanner", func() {
				msg := types.NewMsgQrScan(s1, qr2, nil)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanQRAlreadyOccurred))
			})

			It("should fail if called again by the person scanned", func() {
				msg := types.NewMsgQrScan(s2, qr1, nil)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanQRAlreadyOccurred))
			})
		})

	})
})
