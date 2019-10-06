package key

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/eco/longy/key-service/models"

	"fmt"
)

// GetKeyForEmail retrieves the StoredKey record corresponding to the given
// email address. If an error occurs or no such record can be found an empty
// object may be returned.
//
// The application will crash if unmarshalling fails.
func GetKeyForEmail(db *models.DatabaseContext, email string) (StoredKey) {
	result, err := db.DB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("linkedup-keyservice"),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
	})

	r := StoredKey{}

	if err != nil {
		fmt.Println(err.Error())
		return r
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal StoredKey: %v", err))
	}

	return r
}

// SetStoredKey sets the record associating key material with an email address
// in the application database. A new entry is created if none exists, and the
// existing record is updated if one is already present.
//
// Returns true unless an error occurs.
func SetStoredKey(db *models.DatabaseContext, storedKey *StoredKey) (bool) {
	item, err := dynamodbattribute.MarshalMap(storedKey)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = db.DB.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("linkedup-keyservice"),
		Item: item,
	})

	return err == nil
}
