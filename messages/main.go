package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

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

var slackBotUserOauthToken = os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")
var slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")
var botApi = slack.New(slackBotUserOauthToken)

// var slackOauthToken = os.Getenv("SLACK_OAUTH_ACCESS_TOKEN")
// var regularApi = slack.New(slackOauthToken)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	body := []byte(request.Body)
	log.Printf("Receiving request body %s", body)
	resp := Response{
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200,
	}
	eventsAPIEvent, e := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: slackVerificationToken}),
	)

	if e != nil {
		log.Println("Could not parse Slack event :'(")
		resp.StatusCode = 500
	}
	log.Printf("Processing an event of outer type %s", eventsAPIEvent.Type)

	handleSlackChallenge(eventsAPIEvent, body, resp)

	if eventsAPIEvent.Type == slackevents.CallbackEvent || eventsAPIEvent.Type == slackevents.AppMention {
		innerEvent := eventsAPIEvent.InnerEvent
		log.Printf("Processing an event of inner data %s", innerEvent.Data)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			resp.StatusCode = 200
			log.Printf("Reacting to app mention event from channel %s", ev.Channel)
			_, _, postError := botApi.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if postError != nil {
				resp.StatusCode = 500
				log.Printf("Error while posting message %s", postError)
			}
		case *slackevents.AppHomeOpenedEvent:
			log.Println("Reacting to app home request event")
			resp.StatusCode = 200
			homeViewForUser := slack.HomeTabViewRequest{
				Type:   "home",
				Blocks: apphome.UserHome("Cyril").Blocks,
			}
			homeViewAsJson, _ := json.Marshal(homeViewForUser)
			log.Printf("Sending view %s", homeViewAsJson)
			_, publishViewError := botApi.PublishView(ev.User, homeViewForUser, ev.View.Hash)
			if publishViewError != nil {
				resp.StatusCode = 500
				log.Println(publishViewError)
			}
		}
	}
	return resp, nil
}

// Slack Challenge is used to register the URL in the slack API config interface
// Should only be used once by slack when changing the events URL
func handleSlackChallenge(eventsAPIEvent slackevents.EventsAPIEvent, body []byte, resp Response) {
	if eventsAPIEvent.Type == slackevents.URLVerification {
		buf := new(bytes.Buffer)
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			resp.StatusCode = 500
		}
		resp.Headers["Content-Type"] = "text"
		resp.StatusCode = 200
		buf.Write([]byte(r.Challenge))
		resp.Body = buf.String()
	}
}

func main() {
	lambda.Start(Handler)
}
