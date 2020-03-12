package main

import (
	"dont-slack-evil/notifications"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(notifications.LambdaHandlerLeaderboard)
}