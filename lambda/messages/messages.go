package main

import (
	"dont-slack-evil/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(messages.LambdaHandler)
}
