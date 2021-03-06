package hello

import (
	"context"
	"encoding/json"
	"os"

	"dont-slack-evil/nlp"

	"github.com/aws/aws-lambda-go/events"
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

	sentimentAnalysis, sentimentError := nlp.GetSentiment(message, apiURL, apiKey)
	if sentimentError != nil {
		return Response{StatusCode: 500}, sentimentError
	}
	jsonBody, err := json.Marshal(sentimentAnalysis)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(jsonBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
