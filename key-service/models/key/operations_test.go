package key_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/eco/longy/key-service/models"
	"github.com/eco/longy/key-service/models/key"

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
			Region: aws.String("us-west-1"),
			Endpoint: aws.String("http://localhost:8000"),
		}))

		context = models.NewDatabaseContext(sess)

		_, err := context.DB.CreateTable(&dynamodb.CreateTableInput{
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName: aws.String("linkedup-keyservice"),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("Email"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Email"),
					KeyType: aws.String("HASH"),
				},
			},
		})

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
		}
	})

	It("can save an item to the datastore", func() {
		if !testsEnabled {
			Skip("Use ENABLE_DB_TESTS to turn these on")
		}
		k := key.StoredKey{
			Email: "test@linkedup.sfblockchainweek.io",
			KeyData: []byte{0x0000},
		}

		Expect(key.SetStoredKey(&context, &k)).Should(BeTrue())
	})

	Context("when there is an item in the store", func() {
		var k string
		
		BeforeEach(func() {
			if !testsEnabled {
				return
			}

			k = "test@linkedup.sfblockchainweek.io"

			key.SetStoredKey(&context, &key.StoredKey{
				Email: k,
				KeyData: []byte{0x0000},
			});
		})

		It("can be retrieved", func() {
			if !testsEnabled {
				Skip("Use ENABLE_DB_TESTS to turn these on")
			}
			r := key.GetKeyForEmail(&context, k)

			Expect(r.Email).Should(BeEquivalentTo(k))
		})
	})
})
