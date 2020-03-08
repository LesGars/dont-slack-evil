package notifications

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(ctx context.Context) (Response, error) {
	err := SendNotifications()
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

func main() {
	lambda.Start(LambdaHandler)
}
