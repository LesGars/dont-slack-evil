package main

import (
	"dont-slack-evil/interactive"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(interactive.LambdaHandler)
}
