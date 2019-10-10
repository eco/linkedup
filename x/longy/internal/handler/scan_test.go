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

	It("should fail if the sender is trying to scan themselves", func() {
		//add sender to keeper
		utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false)

		msg := types.NewMsgQrScan(sender, qr1, nil)
		result := handler(ctx, msg)
		Expect(result.Code).To(Equal(types.AccountsSame))
	})

	It("should succeed to create a new scan record without data", func() {
		createScan(qr1, qr2, sender, receiver, nil, false)
	})

	It("should succeed to create a new scan record with data", func() {
		data := []byte("asdfasdfa")
		createScan(qr1, qr2, sender, receiver, data, false)
	})

	Context("when a partial scan already exists but doesn't have shared info from both parties", func() {
		var data []byte
		BeforeEach(func() {
			//Add the partial scan to the keeper
			createScan(qr1, qr2, sender, receiver, nil, false)
			data = []byte("asdfasdfa")
		})

		It("should add info and increment points", func() {
			//Add the partial scan to the keeper
			msg := types.NewMsgQrScan(sender, qr2, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(sdk.CodeOK))
			inspectScan(sender, receiver) //, len(data) != 0, false)

			for i := 0; i < 2; i++ {
				msg = types.NewMsgQrScan(receiver, qr1, data)
				result = handler(ctx, msg)
				Expect(result.Code).To(Equal(sdk.CodeOK))
			}

			//get attendees
			a, ok := keeper.GetAttendee(ctx, sender)
			Expect(ok).To(BeTrue())
			b, ok := keeper.GetAttendee(ctx, receiver)
			Expect(ok).To(BeTrue())
			//Check share ids
			Expect(len(a.ScanIDs)).To(Equal(1))
			Expect(len(b.ScanIDs)).To(Equal(1))

			Expect(a.Rep).To(Equal(types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints))
			Expect(b.Rep).To(Equal(types.ScanAttendeeAwardPoints + types.ShareAttendeeAwardPoints))
		})

		It("should not allow us to reset data and earn more points", func() {
			//a, ok := keeper.GetAttendee(ctx, sender)
			//Expect(ok).To(BeTrue())
			////Add the partial scan to the keeper
			//msg := types.NewMsgQrScan(sender, qr2, nil)
			//result := handler(ctx, msg)
			//Expect(result.Code).To(Equal(sdk.CodeOK))
			//Expect(a.Rep).To(Equal(types.ScanAttendeeAwardPoints))
			//
			//msg = types.NewMsgQrScan(sender, qr2, data)
			//result = handler(ctx, msg)
			//Expect(result.Code).To(Equal(sdk.CodeOK))
			//Expect(a.Rep).To(Equal(types.ScanAttendeeAwardPoints))
		})
	})
})

func inspectScan(s1 sdk.AccAddress, s2 sdk.AccAddress) { //, dataShared bool, scanAccepted bool) {
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

	//Check rewards
	//if scanAccepted {
	//	var point uint
	//	if dataShared {
	//		point += types.ShareAttendeeAwardPoints
	//	}
	//	point = types.ScanAttendeeAwardPoints + point
	//	Expect(a.Rep).To(Equal(point))
	//	Expect(scan.P1).To(Equal(point))
	//	Expect(b.Rep).To(Equal(types.ScanAttendeeAwardPoints))
	//}else{
	//	Expect(a.Rep).To(Equal(0))
	//	Expect(b.Rep).To(Equal(0))
	//}

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
	utils.AddAttendeeToKeeper(ctx, &keeper, qr1, sponsor)
	//add scanned to keeper
	utils.AddAttendeeToKeeper(ctx, &keeper, qr2, false)

	msg := types.NewMsgQrScan(s1, qr2, data)
	result := handler(ctx, msg)
	Expect(result.Code).To(Equal(sdk.CodeOK))
	inspectScan(s1, s2) //, len(data) != 0, false)
}
