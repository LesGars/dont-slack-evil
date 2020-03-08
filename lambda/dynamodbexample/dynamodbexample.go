package main

import (
	dynamodbexample "dont-slack-evil/dynamodbexample"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(dynamodbexample.Handler)
}
