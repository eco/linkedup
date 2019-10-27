package utils_test

import (
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Eventbrite Attendee Genesis Tests", func() {

	It("should set the accounts public address", func() {
		ga := utils.EventbriteAttendee{
			ID:              "123",
			TicketClassName: "Standard",
			Profile:         utils.EventbriteProfile{},
		}
		a := ga.ToGenesisAttendee()
		Expect(a.Sponsor).To(BeFalse())
		Expect(a.Address).To(Equal(util.IDToAddress(ga.ID)))
	})

	It("should not be sponsor when standard ticket type", func() {
		ga := utils.EventbriteAttendee{
			ID:              "123",
			TicketClassName: "Standard",
			Profile:         utils.EventbriteProfile{},
		}
		a := ga.ToGenesisAttendee()
		Expect(a.Sponsor).To(BeFalse())
	})

	It("should be sponsor when sponsor ticket type", func() {
		ga := utils.EventbriteAttendee{
			ID:              "123",
			TicketClassName: "Sponsors",
			Profile:         utils.EventbriteProfile{},
		}
		a := ga.ToGenesisAttendee()
		Expect(a.Sponsor).To(BeTrue())
	})

	It("should be sponsor when cesc speaker ticket type", func() {
		ga := utils.EventbriteAttendee{
			ID:              "123",
			TicketClassName: "CESC Speakers",
			Profile:         utils.EventbriteProfile{},
		}
		a := ga.ToGenesisAttendee()
		Expect(a.Sponsor).To(BeTrue())
	})

	It("should be sponsor when epicenter speaker ticket type", func() {
		ga := utils.EventbriteAttendee{
			ID:              "123",
			TicketClassName: "Epicenter Speakers",
			Profile:         utils.EventbriteProfile{},
		}
		a := ga.ToGenesisAttendee()
		Expect(a.Sponsor).To(BeTrue())
	})
})
