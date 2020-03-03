package main

import (
	"encoding/json"
	"context"
	"fmt"
	"os"

	"dont-slack-evil/nlp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	apiKey := os.Getenv("PD_API_KEY")
	apiURL := os.Getenv("PD_API_URL")
	message := "Go Serverless v1.0! Your function executed successfully!"

	sentimentAnalysis := nlp.GetSentiment(message, apiURL, apiKey)
	jsonBody, err := json.Marshal(map[string]interface{}{
	  "message": sentimentAnalysis.Sentiment,
	  "sentiment": sentimentAnalysis.Message,
	})
	if err != nil {
	  return Response{StatusCode: 404}, err
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(jsonBody),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
