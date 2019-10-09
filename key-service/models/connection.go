package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	infoTableName = "linkedup-keyservice"
	authTableName = "linkedup-keyservice-auth"
)

// DatabaseContext carries context needed to interact with the database.
type DatabaseContext struct {
	DB *dynamodb.DynamoDB
}

// NewDatabaseContext will establish a session with the backend db.
//
// `DatabaseContext` effectively acts as a key-value store for a variety of operations
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

	log.Info("establishing session with dynamo")

	// try create the tables if they haven't already been instantiated
	if err := createTables(context.DB); err != nil {
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

func createTables(db *dynamodb.DynamoDB) error {

	/** create table to store attendee information **/
	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(infoTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
	})

	if err != nil {
		return err
	}

	/** create table to store auth tokens for key recovery **/
	_, err = db.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(authTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
	})

	return err
}

/** Storage **/

// StoreAttendeeInfo -
func (db DatabaseContext) StoreAttendeeInfo(id int, info []byte) bool {
	value := &storedInfo{
		ID:   id,
		Data: info,
	}

	return setInfo(&db, value)
}

// StoreAuthToken -
func (db DatabaseContext) StoreAuthToken(id int, token string) bool {
	auth := &storedAuth{
		ID:        id,
		AuthToken: token,
	}

	return setAuthToken(&db, auth)
}

/** Retrieval **/

// GetAttendeeInfo -
func (db DatabaseContext) GetAttendeeInfo(id int) []byte {
	return getInfoForID(&db, id)
}

// GetAuthToken -
func (db DatabaseContext) GetAuthToken(id int) string {
	value := getAuthTokenForID(&db, id)
	if value == nil {
		return ""
	}

	return value.AuthToken
}
