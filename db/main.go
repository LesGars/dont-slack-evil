package db

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

// Store stores an item in the database
func Store(tableName string, item map[string]string) bool {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("Got error marshalling new item:")
		log.Println(err.Error())
		return false
	}
	putItemInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(putItemInput)
	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())
		return false
	}
	return true
}