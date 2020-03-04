package db

import (
	"strings"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


// CreateDBIfNotCreated creates DynamoDB table if it doesn't exist
func CreateDBIfNotCreated(tableName string) bool {
	// Create DynamoDB client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName: aws.String(tableName),
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}
	_, createTableErr := svc.CreateTable(createTableInput)
	if createTableErr != nil {
		if (!strings.Contains(createTableErr.Error(), "Table already exists")) {
			log.Println(createTableErr.Error())
			return false;
		}
	} else {
		log.Println("Created the table", tableName)
	}
	return true
}
