package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "key-storage")

// getKeyForID retrieves the information record corresponding to the given
// email address. If a record does not exist, nil bytes will be returned
//
// The application will crash if unmarshalling fails.
func getInfoForID(db *DatabaseContext, id int) ([]byte, error) {
	result, err := db.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(infoTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(idToString(id)),
			},
		},
	})
	if err != nil {
		log.WithError(err).WithField("id", id).Info("failed info retrieval")
		return nil, err
	} else if result == nil || result.Item == nil || len(result.Item) == 0 {
		// item not found
		return nil, nil
	}

	var r storedInfo
	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal storedInfo: %v", err))
	}

	return r.Data, nil
}

// getVerificationTokenForID retrieves the auth record corresponding to the given
// email address. `nil` will be returned for entries that do not exist
//
// The application will crash if unmarshalling fails.
func getVerificationTokenForID(db *DatabaseContext, id int) (string, error) {
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
		return "", err
	} else if result == nil || result.Item == nil || len(result.Item) == 0 {
		// item not found
		return "", nil
	}

	var r storedAuth
	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal StoredAuth: %s", err))
	}

	return r.AuthToken, nil
}

func getEmailForID(db *DatabaseContext, id int) *storeEmail {
	result, err := db.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(emailTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(idToString(id)),
			},
		},
	})

	if err != nil {
		log.WithError(err).WithField("id", id).Info("no email for id stored")
		return nil
	}

	var e storeEmail
	err = dynamodbattribute.UnmarshalMap(result.Item, &e)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal storeEmail: %s", err))
	}

	return &e
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
		TableName: aws.String(infoTableName),
		Item:      item,
	})

	if err != nil {
		log.WithError(err).Error("failed key storage")
		return false
	}

	return true
}

// setVerificationTokenForEmail sets the record associating key material with an email address
// in the application database. A new entry is created if none exists, and the
// existing record is updated if one is already present.
//
// Returns true unless an error occurs.
func setVerificationToken(db *DatabaseContext, key *storedAuth) bool {
	item, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		panic(err)
	}

	_, err = db.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(authTableName),
		Item:      item,
	})
	if err != nil {
		log.WithError(err).Error("failed auth storage")
		return false
	}

	return true
}

func setEmail(db *DatabaseContext, email *storeEmail) bool {
	item, err := dynamodbattribute.MarshalMap(email)
	if err != nil {
		panic(err)
	}

	_, err = db.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(emailTableName),
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
