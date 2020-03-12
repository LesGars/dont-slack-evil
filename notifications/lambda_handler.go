package notifications

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// LambdaHandler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(ctx context.Context) (Response, error) {
	_, err := SendNotifications()
	if err != nil {
		log.Println(err)
	  return Response{StatusCode: 500}, err
	}

	return Response{
		StatusCode:      204,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type":           "text",
		},
	}, nil
}
