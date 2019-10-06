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
				info, err := types.NewInfo(sender, receiver, data)
				Expect(err).To(BeNil())
				keeper.SetInfo(ctx, &info)

				msg := types.NewMsgInfo(sender, receiver, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.InfoAlreadyExists))

			})

			It("should fail when there is no scan between the sender and receiver ", func() {
				msg := types.NewMsgInfo(sender, receiver, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.ScanNotFound))
			})

			It("should fail when the scan between sender and receiver is not complete ", func() {
				scan, err := types.NewScan(sender, receiver)
				Expect(err).To(BeNil())
				keeper.SetScan(ctx, &scan)

				msg := types.NewMsgInfo(sender, receiver, data)
				result := handler(ctx, msg)
				Expect(result.Code).To(Equal(types.InvalidShareForScan))
			})

			Context("when scan between sender and receiver is complete", func() {
				var senderRep uint
				var recIdsLen int
				BeforeEach(func() {
					//store the init rep
					s, errAtt := keeper.GetAttendee(ctx, sender)
					Expect(errAtt).To(BeNil())
					senderRep = s.Rep

					//store the init ids array length
					s, errAtt = keeper.GetAttendee(ctx, receiver)
					Expect(errAtt).To(BeNil())
					recIdsLen = len(s.InfoIDs)

					scan, err := types.NewScan(sender, receiver)
					Expect(err).To(BeNil())
					scan.Complete = true
					keeper.SetScan(ctx, &scan)
				})
				FIt("should succeed when participants are both regular attendees", func() {
					msg := types.NewMsgInfo(sender, receiver, data)
					result := handler(ctx, msg)
					Expect(result.Code).To(Equal(sdk.CodeOK))

					//receiver should have info id
					s2, err := keeper.GetAttendee(ctx, receiver)
					Expect(err).To(BeNil())
					Expect(len(s2.InfoIDs)).To(Equal(recIdsLen + 1))
					Expect(Contains(s2.InfoIDs))
				})

				It("should succeed when one participant is a sponsor attendee", func() {
				})
			})

		})
	})
})

func Contains(vals []string, v string) (contains bool) {
	for _, a := range vals {
		if a == v {
			return true
		}
	}
	return false
}
