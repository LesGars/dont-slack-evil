package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"dont-slack-evil/apphome"
	dsedb "dont-slack-evil/db"
	"dont-slack-evil/nlp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fatih/structs"
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

var slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")

// If you need to send a message from the app's "bot user", use this bot client
var slackBotUserOauthToken = os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")
var slackBotUserApiClient = slack.New(slackBotUserOauthToken)

// If you need anything else, use this client instead
// var slackOauthToken = os.Getenv("SLACK_OAUTH_ACCESS_TOKEN")
// var slackRegularApiClient = slack.New(slackOauthToken)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	structs.DefaultTagName = "json" // https://github.com/fatih/structs/issues/25
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

	handleSlackChallenge(eventsAPIEvent, body, &resp)

	if eventsAPIEvent.Type == slackevents.CallbackEvent || eventsAPIEvent.Type == slackevents.AppMention {
		innerEvent := eventsAPIEvent.InnerEvent
		log.Printf("Processing an event of inner data %s", innerEvent.Data)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			message := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
			storeMessage(message, &resp)
			getSentiment(message, &resp)
		case *slackevents.AppMentionEvent:
			resp.StatusCode = 200
			log.Printf("Reacting to app mention event from channel %s", ev.Channel)
			_, _, postError := slackBotUserApiClient.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if postError != nil {
				resp.StatusCode = 500
				log.Printf("Error while posting message %s", postError)
			}
		case *slackevents.AppHomeOpenedEvent:
			log.Println("Reacting to app home request event")
			resp.StatusCode = 200
			homeViewForUser := slack.HomeTabViewRequest{
				Type:   "home",
				Blocks: apphome.UserHome(ev.User).Blocks,
			}
			homeViewAsJson, _ := json.Marshal(homeViewForUser)
			log.Printf("Sending view %s", homeViewAsJson)
			_, publishViewError := slackBotUserApiClient.PublishView(ev.User, homeViewForUser, ev.View.Hash)
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
func handleSlackChallenge(eventsAPIEvent slackevents.EventsAPIEvent, body []byte, resp *Response) {
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

func storeMessage(message *slackevents.MessageEvent, resp *Response) {
	// Create DB
	tableName := os.Getenv("DYNAMODB_TABLE")
	dbError := dsedb.CreateDBIfNotCreated(tableName)
	if dbError {
		resp.StatusCode = 500
	}

	// Save in DB
	messageBytes, _ := json.Marshal(message)
	dbItem := dsedb.Message{
		UserId:         message.User,
		Text:           message.Text,
		CreatedAt:      message.EventTimeStamp.String(),
		SlackMessageId: message.EventTimeStamp.String(),
		SlackThreadId:  message.ThreadTimeStamp,
	}
	json.Unmarshal(messageBytes, &dbItem)
	log.Println(structs.Map(&dbItem))

	dbResult := dsedb.Store(tableName, structs.Map(&dbItem))
	if !dbResult {
		log.Println("Could not store message in DB")
	} else {
		log.Println("Message was stored successfully")
	}
}

func getSentiment(message *slackevents.MessageEvent, resp *Response) {
	tableName := os.Getenv("DYNAMODB_TABLE")
	apiKey := os.Getenv("PD_API_KEY")
	apiURL := os.Getenv("PD_API_URL")
	text := message.Text
	sentimentAnalysis, sentimentError := nlp.GetSentiment(text, apiURL, apiKey)
	if sentimentError != nil {
		log.Println("Could not analyze message")
		resp.StatusCode = 500
	}
	dbResult := dsedb.Update(tableName, message.EventTimeStamp.String(), sentimentAnalysis.Sentiment)
	if !dbResult {
		log.Println("Could not update message with sentiment")
	} else {
		log.Println("Message was updated successfully with sentiment")
	}
}

func main() {
	lambda.Start(Handler)
}
