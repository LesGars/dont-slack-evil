package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "dont-slack-evil/messages/service"
	// "github.com/fatih/structs"

)

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	body := []byte(request.Body)
	query := request.QueryStringParameters
	log.Printf("Receiving request body %s", body)
	log.Printf("Receiving request query string %s", query)
	redirectUrl := "http://lol.com"
	resp := Response{
		StatusCode: 200,
		Body: redirectUrl,
		// StatusCode: 302,
		// Headers: map[string]string{
		// 	"Location": redirectUrl,
		// },
		// Body: "",
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
