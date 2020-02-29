package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function for the incoming webhook endpoint
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	webhookURL := os.Getenv("SLACK_INCOMING_WEBHOOK_URL")

	reqBody, err := json.Marshal(map[string]string{
		"text": "Hello world from the incoming-webhook lambda",
	})
	if err != nil {
			fmt.Println(err.Error())
	}
	resp, err := http.Post(webhookURL,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return events.APIGatewayProxyResponse {
		Body: string(body),
		StatusCode: 200,
		Headers: map[string]string{
			"content-type": "text/plain",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
