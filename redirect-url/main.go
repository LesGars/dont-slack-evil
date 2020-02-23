package main

import (
	"fmt"
	"os"
	// "net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Item lol
type Item struct {
	Year   int
	Title  string
	Plot   string
	Rating float64
}

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
	dynamoTable := os.Getenv("DYNAMODB_TABLE")
	fmt.Println(dynamoTable)
	url := fmt.Sprintf(
		"https://slack.com/api/oauth.access?client_id=%s&client_secret=%s&code=%s&state=%s",
		clientID,
		clientSecret,
		code,
		state,
	)
	fmt.Println(url)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_table.session]

	// snippet-start:[dynamodb.go.create_table.call]
	// Create table Movies
	tableName := dynamoTable 

	// input := &dynamodb.CreateTableInput{
	// 		AttributeDefinitions: []*dynamodb.AttributeDefinition{
	// 				{
	// 						AttributeName: aws.String("Year"),
	// 						AttributeType: aws.String("N"),
	// 				},
	// 				{
	// 						AttributeName: aws.String("Title"),
	// 						AttributeType: aws.String("S"),
	// 				},
	// 		},
	// 		KeySchema: []*dynamodb.KeySchemaElement{
	// 				{
	// 						AttributeName: aws.String("Year"),
	// 						KeyType:       aws.String("HASH"),
	// 				},
	// 				{
	// 						AttributeName: aws.String("Title"),
	// 						KeyType:       aws.String("RANGE"),
	// 				},
	// 		},
	// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
	// 				ReadCapacityUnits:  aws.Int64(10),
	// 				WriteCapacityUnits: aws.Int64(10),
	// 		},
	// 		TableName: aws.String(tableName),
	// }

	// _, err := svc.CreateTable(input)
	// if err != nil {
	// 		fmt.Println("Got error calling CreateTable:")
	// 		fmt.Println(err.Error())
	// 		os.Exit(1)
	// }

	// fmt.Println("Created the table", tableName)





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

	item := Item{
			Year:   2015,
			Title:  "The Big New Movie",
			Plot:   "Nothing happens at all.",
			Rating: 0.0,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
			fmt.Println("Got error marshalling new movie item:")
			fmt.Println(err.Error())
			os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
	}

	year := strconv.Itoa(item.Year)

	fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse {
		Body: "ok",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
