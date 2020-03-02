package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	var buf bytes.Buffer

	apiKey := os.Getenv("PD_API_KEY")
	apiURL := os.Getenv("PD_API_URL")
	message := "Go Serverless v1.0! Your function executed successfully!"

	// Get sentiment of message
	form := url.Values{}
	form.Add("text", message)
	form.Add("api_key", apiKey)
	resp, _ := http.Post(
		apiURL + "/v4/sentiment",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	responseBody, err := json.Marshal(map[string]interface{}{
		"message": message,
		"sentiment": string(body),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, responseBody)

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(responseBody),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
