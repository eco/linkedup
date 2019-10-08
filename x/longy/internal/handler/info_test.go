package handler_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = PDescribe("Info Handler Tests", func() {
	var sender, receiver sdk.AccAddress
	var data = []byte{1, 2, 3, 2, 1}
	const (
		qr1 = "1234"
		qr2 = "asdf"
	)
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
			utils.AddAttendeeToKeeper(ctx, &keeper, qr1, false)
		})

		It("should fail when receiver attendee doesn't exist", func() {
			msg := types.NewMsgInfo(sender, receiver, data)
			result := handler(ctx, msg)
			Expect(result.Code).To(Equal(types.AttendeeNotFound))
		})

		Context("when both sender and receiver exist", func() {
			BeforeEach(func() {
				utils.AddAttendeeToKeeper(ctx, &keeper, qr2, false)
			})

			It("should fail when there is already a shared info between sender and receiver", func() {
				//info, err := types.NewInfo(sender, receiver, data)
				//Expect(err).To(BeNil())
				////keeper.SetInfo(ctx, &info)
				//
				//msg := types.NewMsgInfo(sender, receiver, data)
				//result := handler(ctx, msg)
				//Expect(result.Code).To(Equal(types.InfoAlreadyExists))

			})

			It("should fail when there is no scan between the sender and receiver ", func() {
				msg := types.NewMsgInfo(sender, receiver, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanNotFound))
			})

			It("should fail when the scan between sender and receiver is not complete ", func() {
				scan, err := types.NewScan(sender, receiver, nil, nil)
				Expect(err).To(BeNil())
				keeper.SetScan(ctx, &scan)

				msg := types.NewMsgInfo(sender, receiver, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.InvalidShareForScan))
			})

			Context("when scan between sender and receiver is complete", func() {
				var senderRep, receiverRep uint
				var recIdsLen int
				BeforeEach(func() {
					//store the init rep
					s, ok := keeper.GetAttendee(ctx, sender)
					Expect(ok).To(BeTrue())
					senderRep = s.Rep

					//store the init ids array length
					s, ok = keeper.GetAttendee(ctx, receiver)
					Expect(ok).To(BeTrue())
					receiverRep = s.Rep
					recIdsLen = len(s.ScanIDs)

					scan, err := types.NewScan(sender, receiver, nil, nil)
					Expect(err).To(BeNil())
					scan.Complete = true
					keeper.SetScan(ctx, &scan)
				})

				It("should succeed when participants are both regular attendees", func() {
					msg := types.NewMsgInfo(sender, receiver, data)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))

					//receiver should have info id
					s2, ok := keeper.GetAttendee(ctx, receiver)
					Expect(ok).To(BeTrue())
					Expect(len(s2.ScanIDs)).To(Equal(recIdsLen + 1))
					id, e := types.GenInfoID(sender, receiver)
					Expect(e).To(BeNil())
					Expect(Contains(s2.ScanIDs, string(id))).To(BeTrue())

					s, ok := keeper.GetAttendee(ctx, sender)
					Expect(ok).To(BeTrue())
					Expect(s.Rep).To(Equal(senderRep + types.ShareAttendeeAwardPoints))

					p, ok := keeper.GetAttendee(ctx, receiver)
					Expect(ok).To(BeTrue())
					Expect(p.Rep).To(Equal(receiverRep)) //shouldn't change
				})

				It("should succeed when one participant is a sponsor attendee", func() {
					utils.AddAttendeeToKeeper(ctx, &keeper, qr2, true)

					msg := types.NewMsgInfo(sender, receiver, data)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))

					//receiver should have info id
					s2, ok := keeper.GetAttendee(ctx, receiver)
					Expect(ok).To(BeTrue())
					Expect(len(s2.ScanIDs)).To(Equal(recIdsLen + 1))
					id, e := types.GenInfoID(sender, receiver)
					Expect(e).To(BeNil())
					Expect(Contains(s2.ScanIDs, string(id))).To(BeTrue())

					s, ok := keeper.GetAttendee(ctx, sender)
					Expect(ok).To(BeTrue())
					Expect(s.Rep).To(Equal(senderRep + types.ShareSponsorAwardPoints))

					p, ok := keeper.GetAttendee(ctx, receiver)
					Expect(ok).To(BeTrue())
					Expect(p.Rep).To(Equal(receiverRep)) //shouldn't change
				})
			})

		})
	})
})

//Contains checks to see if a value is in the array
func Contains(vals []string, v string) (contains bool) {
	for _, a := range vals {
		if a == v {
			return true
		}
	}
	return false
}
