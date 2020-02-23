package main

import (
	"fmt"
	"os"
	"strings"
	// "net/http"
	// "strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	code := request.QueryStringParameters["code"]
	state := request.QueryStringParameters["state"]
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	url := fmt.Sprintf(
		"https://slack.com/api/oauth.access?client_id=%s&client_secret=%s&code=%s&state=%s",
		clientID,
		clientSecret,
		code,
		state,
	)
	fmt.Println(url)

	// Create DynamoDB client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	// Create DynamoDB table if it doesn't exist
	tableName := os.Getenv("DYNAMODB_TABLE")
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
	_, err := svc.CreateTable(createTableInput)
	if err != nil {
		if (!strings.Contains(err.Error(), "Table already exists")) {
			fmt.Println("Got error calling CreateTable:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("Created the table", tableName)
	}

	// Add item to DynamoDB table
	item := map[string]string{
		"id": "bonjour",
	}
	fmt.Println(item)
	av, err := dynamodbattribute.MarshalMap(item)
	fmt.Println(av)
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
	} else {
		fmt.Println("Successfully added to" + tableName)
	}


	// resp, err := http.Get(url)
	// fmt.Println(resp)
	// fmt.Println(err)

  // const params = {
  //   TableName: process.env.DYNAMODB_TABLE,
  //   Item: {
  //     id: item.team_id,
  //     ...item,
  //   },
	// };


	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse {
		Body: "ok",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
