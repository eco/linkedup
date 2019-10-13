package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/sirupsen/logrus"

	"fmt"
)

var log = logrus.WithField("module", "key-storage")

// getKeyForID retrieves the information record corresponding to the given
// email address. If an error occurs or no such record can be found an empty
// object may be returned.
//
// The application will crash if unmarshalling fails.
func getInfoForID(db *DatabaseContext, id int) []byte {
	result, err := db.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(infoTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(idToString(id)),
			},
		},
	})

	if err != nil {
		log.WithError(err).WithField("id", id).Info("failed key retrieval")
		return nil
	}

	var r storedInfo
	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal storedInfo: %v", err))
	}

	return r.Data
}

// getAuthTokenForID retrieves the auth record corresponding to the given
// email address. If an error occurs or no such record can be found an empty
// object may be returned.
//
// The application will crash if unmarshalling fails.
func getAuthTokenForID(db *DatabaseContext, id int) *storedAuth {
	result, err := db.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(authTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(idToString(id)),
			},
		},
	})

	if err != nil {
		log.WithError(err).WithField("id", id).Info("failed auth retrieval")
		return nil
	}

	var r storedAuth
	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal StoredAuth: %s", err))
	}

	return &r
}

// setKey sets the record associating key material with an email address
// in the application database. A new entry is created if none exists, and the
// existing record is updated if one is already present.
//
// Returns true unless an error occurs.
func setInfo(db *DatabaseContext, key *storedInfo) bool {
	item, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		panic(err)
	}

	_, err = db.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("linkedup-keyservice"),
		Item:      item,
	})

	if err != nil {
		log.WithError(err).Error("failed key storage")
		return false
	}

	return true
}

// setAuthTokenForEmail sets the record associating key material with an email address
// in the application database. A new entry is created if none exists, and the
// existing record is updated if one is already present.
//
// Returns true unless an error occurs.
func setAuthToken(db *DatabaseContext, key *storedAuth) bool {
	item, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		panic(err)
	}

	_, err = db.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("linkedup-keyservice-auth"),
		Item:      item,
	})

	if err != nil {
		log.WithError(err).Error("failed auth storage")
		return false
	}

	return true
}

/** Helpers **/
func idToString(i int) string {
	return fmt.Sprintf("%d", i)
}
