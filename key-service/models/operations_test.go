package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/eco/longy/key-service/models"

	"fmt"
	"os"
)

var _ = Describe("StoredKey Operations", func() {
	var sess *session.Session
	var context models.DatabaseContext
	var testsEnabled bool

	BeforeSuite(func() {
		if os.Getenv("ENABLE_DB_TESTS") != "true" {
			testsEnabled = false
		} else {
			testsEnabled = true
		}

		sess = session.Must(session.NewSession(&aws.Config{
			Region:   aws.String("us-west-1"),
			Endpoint: aws.String("http://localhost:8000"),
		}))

		var err error
		context, err = models.NewDatabaseContextWithCfg(sess)
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	It("can save an item to the datastore", func() {
		if !testsEnabled {
			Skip("Use ENABLE_DB_TESTS to turn these on")
		}

		Expect(context.StoreKey("howdy", "hi")).Should(BeTrue())
	})

	Context("when there is an item in the store", func() {
		var email string
		var data = "hi"

		BeforeEach(func() {
			if !testsEnabled {
				return
			}

			email = "test@linkedup.sfblockchainweek.io"
			context.StoreKey(email, "hi")
		})

		It("can be retrieved", func() {
			if !testsEnabled {
				Skip("Use ENABLE_DB_TESTS to turn these on")
			}
			res := context.GetKey(email)

			Expect(res).Should(BeEquivalentTo(data))
		})
	})
})
