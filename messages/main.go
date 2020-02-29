package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	//"dont-slack-evil/apphome"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var slackOauthToken = os.Getenv("SLACK_OAUTH_TOKEN")
var api = slack.New(slackOauthToken)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, body []byte) (Response, error) {
	buf := new(bytes.Buffer)
	resp := Response{
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: slackOauthToken}))
	if e != nil {
		resp.StatusCode = 500
	}

	if eventsAPIEvent.Type == slackevents.Message {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
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
			//TODO Use the slack.HomeTabViewRequest struct correctly, and find the hash

			//api.PublishView(ev.User, slack.HomeTabViewRequest{"home", apphome.UserHome("Cyril")}, "")
		}
	}

	return resp, nil

}

func main() {
	lambda.Start(Handler)
}
