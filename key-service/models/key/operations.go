package key

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/eco/longy/key-service/models"

	"fmt"
)

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

func SetStoredKey(db *models.DatabaseContext, storedKey *StoredKey) (bool) {
	item, err := dynamodbattribute.MarshalMap(storedKey)
	_, err = db.DB.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("linkedup-keyservice"),
		Item: item,
	})

	return err == nil
}
