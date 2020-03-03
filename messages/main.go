package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"dont-slack-evil/messages/service"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	body := []byte(request.Body)
	log.Printf("Receiving request body %s", body)
	resp := Response{
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200,
	}
	challengeResponse, err := service.HandleEvent(body)
	if err != nil {
		resp.StatusCode = 500
	} else {
		resp.StatusCode = 200
		resp.Body = challengeResponse
		resp.Headers["Content-Type"] = "text"
	}
	
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
