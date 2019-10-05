package models

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DatabaseContext struct {
	DB   *dynamodb.DynamoDB
}

func NewDatabaseContext(sess *session.Session) (context DatabaseContext) {
	context = DatabaseContext{}

	context.DB = dynamodb.New(sess)

	return
}
