package models

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DatabaseContext carries context needed to interact with the database.
type DatabaseContext struct {
	DB   *dynamodb.DynamoDB
}

// NewDatabaseContext constructs a new DatabaseContext, using the given AWS
// session handle.
func NewDatabaseContext(cfg client.ConfigProvider) (context DatabaseContext) {
	context = DatabaseContext{}

	context.DB = dynamodb.New(cfg)

	return
}
