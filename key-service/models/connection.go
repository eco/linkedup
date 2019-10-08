package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DatabaseContext carries context needed to interact with the database.
type DatabaseContext struct {
	DB *dynamodb.DynamoDB
}

// NewDatabaseContext will establish a session with the backend db
func NewDatabaseContext(region string, endpoint string) (DatabaseContext, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}))

	return NewDatabaseContextWithCfg(sess)
}

// NewDatabaseContextWithCfg constructs a new DatabaseContext, using the given AWS
// session handle.
func NewDatabaseContextWithCfg(cfg client.ConfigProvider) (DatabaseContext, error) {
	context := DatabaseContext{}
	context.DB = dynamodb.New(cfg)

	log.Info("etablishing session with dynamo")

	// create the tables
	_, err := context.DB.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String("linkedup-keyservice"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       aws.String("HASH"),
			},
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
				log.Info("dynamo table already created")
				return context, nil
			}
		} else {
			return context, err
		}
	}

	return context, nil
}

// StoreKey -
func (db DatabaseContext) StoreKey(email string, key string) bool {
	value := &storedKey{
		Email:   email,
		KeyData: []byte(key),
	}

	return setStoredKey(&db, value)
}

// GetKey -
func (db DatabaseContext) GetKey(email string) string {
	value := getKeyForEmail(&db, email)
	return string(value.KeyData)
}
