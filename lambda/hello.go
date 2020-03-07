package main

import (
	"dont-slack-evil/hello"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(hello.Handler)
}
