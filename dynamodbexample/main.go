package dynamodbexample

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Create DynamoDB client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	// Create DynamoDB table if it doesn't exist
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "example"
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
		TableName:   aws.String(tableName),
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}
	_, err := svc.CreateTable(createTableInput)
	if err != nil {
		if !strings.Contains(err.Error(), "Table already exists") {
			fmt.Println("Got error calling CreateTable:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("Created the table", tableName)
	}

	// Add item to DynamoDB table
	rand.Seed(time.Now().UnixNano())
	randomString := strconv.Itoa(rand.Int())
	item := map[string]string{
		"id": randomString,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	putItemInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(putItemInput)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Retrieve item from DynamoDB table
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(randomString),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%s", result),
		StatusCode: 200,
		Headers: map[string]string{
			"content-type": "text/plain",
		},
	}, nil
}
