package utils_test

import (
	"fmt"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Longy Module Test Suite")
}

var _ = Describe("Claim Id Tests", func() {
	var key string
	BeforeEach(func() {
		key = os.Getenv(utils.EventbriteEnvKey)
		err := os.Unsetenv(utils.EventbriteEnvKey)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		if key != "" {
			_ = os.Setenv(utils.EventbriteEnvKey, key)
		} else {
			err := os.Unsetenv(utils.EventbriteEnvKey)
			Expect(err).To(BeNil())
		}
	})

	It("should fail when environment variable for auth not set", func() {
		_, err := utils.GetAttendees()
		Expect(err).To(Not(BeNil()))
	})

	It("should fail when auth header is invalid", func() {
		errOS := os.Setenv(utils.EventbriteEnvKey, "afakekeyandstuff")
		Expect(errOS).To(BeNil())
		_, err := utils.GetAttendees()
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.NetworkResponseError))
	})

	It("should succeed to get all the attendees", func() {
		sysErr := os.Setenv(utils.EventbriteEnvKey, key)
		Expect(sysErr).To(BeNil())

		_, err := utils.GetAttendees()
		Expect(err).To(BeNil())
	})

	It("should have no duplicate attendees in the result on success", func() {
		sysErr := os.Setenv(utils.EventbriteEnvKey, key)
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
