package main

import (
	"fmt"
	"os"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	resp, err := http.Get(url)
	fmt.Println(resp)
	fmt.Println(err)

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse {
		Body: "ok",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
