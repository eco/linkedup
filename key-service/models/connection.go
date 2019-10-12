package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
	"fmt"
)

const (
	infoTableName = "linkedup-keyservice"
	authTableName = "linkedup-keyservice-auth"
)

// DatabaseContext carries context needed to interact with the database.
type DatabaseContext struct {
	db *dynamodb.DynamoDB
	s3 *s3.S3
	contentBucket string
}

// NewDatabaseContext will establish a session with the backend db.
//
// `DatabaseContext` effectively acts as a key-value store for a variety of operations
func NewDatabaseContext(localstack bool, contentBucket string) (DatabaseContext, error) {
	return NewDatabaseContextWithCfg(
		session.Must(session.NewSession()),
		localstack,
		contentBucket,
	)
}

// NewDatabaseContextWithCfg constructs a new DatabaseContext, using the given AWS
// session handle.
func NewDatabaseContextWithCfg(cfg client.ConfigProvider, localstack bool, bucket string) (DatabaseContext, error) {
	context := DatabaseContext{}
	if localstack {
		context.db = dynamodb.New(
			cfg,
			&aws.Config{
				Endpoint: aws.String("http://localstack:4569"),
			},
		)
		context.s3 = s3.New(
			cfg,
			&aws.Config{
				Endpoint: aws.String("http://localstack:4572"),
			},
		)
	} else {
		context.db = dynamodb.New(cfg)
		context.s3 = s3.New(cfg)
	}
	context.contentBucket = bucket

	log.Info("establishing session with dynamo")

	// try create the tables if they haven't already been instantiated
	if err := createTables(context.db); err != nil {
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

// GetImageUploadURL get a URL that an image can be uploaded to
func (db DatabaseContext) GetImageUploadURL(key string) (string, error) {
	req, _ := db.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(db.contentBucket),
		Key:    aws.String(fmt.Sprintf("avatars/%s", key)),
	})

	result, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", err
	}

	return result, nil
}
