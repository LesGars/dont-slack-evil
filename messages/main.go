package main

import (
	"bytes"
	"encoding/json"
	"os"
	"fmt"

	"dont-slack-evil/apphome"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

var slackOauthToken = os.Getenv("SLACK_OAUTH_ACCESS_TOKEN")
var slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")
var api = slack.New(slackOauthToken)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	buf := new(bytes.Buffer)
	body := []byte(request.Body)
	resp := Response{
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: slackVerificationToken}))

	if e != nil {
		fmt.Println(e)
		resp.StatusCode = 500
	}

	if eventsAPIEvent.Type == slackevents.Message {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			fmt.Println(err)
			resp.StatusCode = 500
		}
		resp.Headers["Content-Type"] = "text"
		buf.Write([]byte(r.Challenge))
		resp.Body = buf.String()
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		case *slackevents.AppHomeOpenedEvent:
			api.PublishView(
				ev.User,
				slack.HomeTabViewRequest{"home", apphome.UserHome("Cyril").Blocks, "", "", ""},
				ev.View.Hash)
		}
	}

	return resp, nil

}

func main() {
	lambda.Start(Handler)
}
