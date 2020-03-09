package interactive

import (
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/fatih/structs"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(request Request) (Response, error) {
	structs.DefaultTagName = "json" // https://github.com/fatih/structs/issues/25
	params, err := url.ParseQuery(request.Body)
	if err != nil {
		log.Println(err)
	}
	payload := params.Get("payload")

	log.Printf("Receiving request body.payload %s", payload)
	resp := Response{
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200,
	}
	handlerErr := SlackHandler(payload, resp)
	if handlerErr != nil {
		resp.StatusCode = 500
	}

	return resp, nil
}
