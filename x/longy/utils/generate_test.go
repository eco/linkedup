package utils_test

import (
	"fmt"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate Attendee Genesis Tests", func() {
	var (
		key      string
		event    string
		isEnvSet bool
	)

	BeforeEach(func() {
		key = os.Getenv(utils.EventbriteAuthEnvKey)
		event = os.Getenv(utils.EventbriteEventEnvKey)

		isEnvSet = event != "" && key != ""
	})

	AfterEach(func() {
		if isEnvSet {
			_ = os.Setenv(utils.EventbriteAuthEnvKey, key)
			_ = os.Setenv(utils.EventbriteEventEnvKey, event)
		} else {
			err := os.Unsetenv(utils.EventbriteAuthEnvKey)
			Expect(err).To(BeNil())
		}
	})

	It("should fail when environment variable for auth not set", func() {
		_ = os.Unsetenv(utils.EventbriteAuthEnvKey)
		_, err := utils.GetAttendees()
		Expect(err).To(Not(BeNil()))
	})

	It("should fail when auth header is invalid", func() {
		if !isEnvSet {
			Skip("set EVENTBRITE_AUTH and EVENTBRITE_EVENT to run tests that access Eventbrite data")
		}

		errOS := os.Setenv(utils.EventbriteAuthEnvKey, "afakekeyandstuff")
		Expect(errOS).To(BeNil())
		_, err := utils.GetAttendees()
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.NetworkResponseError))
	})

	It("should succeed to get all the attendees", func() {
		if !isEnvSet {
			Skip("set EVENTBRITE_AUTH and EVENTBRITE_EVENT to run tests that access Eventbrite data")
		}
		sysErr := os.Setenv(utils.EventbriteAuthEnvKey, key)
		Expect(sysErr).To(BeNil())

		_, err := utils.GetAttendees()
		Expect(err).To(BeNil())
	})

	It("should have no duplicate attendees in the result on success", func() {
		if !isEnvSet {
			Skip("set EVENTBRITE_AUTH and EVENTBRITE_EVENT to run tests that access Eventbrite data")
		}
		sysErr := os.Setenv(utils.EventbriteAuthEnvKey, key)
		Expect(sysErr).To(BeNil())

		ga, err := utils.GetAttendees()
		Expect(err).To(BeNil())

		dup := make(map[string]bool, len(ga))
		for i, a := range ga {
			v := dup[a.ID]
			if v {
				msg := fmt.Sprintf("there is a duplicate attendee in the result with"+
					" id : %s, position : %d", a.ID, i)
				Fail(msg)
			}
			dup[a.ID] = true
		}
	})
})
