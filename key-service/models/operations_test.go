package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

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

		context, err = models.NewDatabaseContextWithCfg(sess)
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	It("can save an item to the datastore", func() {
		if !testsEnabled {
			Skip("Use ENABLE_DB_TESTS to turn these on")
		}
		k := storedKey{
			Email:   "test@linkedup.sfblockchainweek.io",
			KeyData: []byte{0x0000},
		}

		Expect(setStoredKey(&context, &k)).Should(BeTrue())
	})

	Context("when there is an item in the store", func() {
		var k string

		BeforeEach(func() {
			if !testsEnabled {
				return
			}

			k = "test@linkedup.sfblockchainweek.io"

			setStoredKey(&context, &storedKey{
				Email:   k,
				KeyData: []byte{0x0000},
			})
		})

		It("can be retrieved", func() {
			if !testsEnabled {
				Skip("Use ENABLE_DB_TESTS to turn these on")
			}
			r := getKeyForEmail(&context, k)

			Expect(r.Email).Should(BeEquivalentTo(k))
		})
	})
})
