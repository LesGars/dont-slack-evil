package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"dont-slack-evil/apphome"
	dsedb "dont-slack-evil/db"
	"dont-slack-evil/nlp"
	"github.com/fatih/structs"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)


var slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")

// If you need to send a message from the app's "bot user", use this bot client
var slackBotUserOauthToken = os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")
var slackBotUserApiClient = slack.New(slackBotUserOauthToken)

// If you need anything else, use this client instead
// var slackOauthToken = os.Getenv("SLACK_OAUTH_ACCESS_TOKEN")
// var slackRegularApiClient = slack.New(slackOauthToken)


// ParseEvent is the assignation of slackevents.ParseEvent to a variable,
// in order to make it mockable
var parseEvent = slackevents.ParseEvent

// ParseEvent is the assignation of slackBotUserApiClient.PostMessage to a variable,
// in order to make it mockable
var postMessage = slackBotUserApiClient.PostMessage

// PublishView is the assignation of slackBotUserApiClient.PublishView to a variable,
// in order to make it mockable
var publishView = slackBotUserApiClient.PublishView 

// HandleEvent uses Slack's Event API to respond to an event emitted by our application
func HandleEvent(body []byte) (string, error) {
	var challengeResponse string
	eventsAPIEvent, e := parseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: slackVerificationToken}),
	)

	if e != nil {
		const message = "Could not parse Slack event :'("
		log.Println(message)
		return "", errors.New(message)
	}
	log.Printf("Processing an event of outer type %s", eventsAPIEvent.Type)

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var challengeError error
		challengeResponse, challengeError = handleSlackChallenge(eventsAPIEvent, body)
		if challengeError != nil {
			return "", challengeError
		}
	}
		
	if eventsAPIEvent.Type == slackevents.CallbackEvent || eventsAPIEvent.Type == slackevents.AppMention {
		innerEvent := eventsAPIEvent.InnerEvent
		log.Printf("Processing an event of inner data %s", innerEvent.Data)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			log.Printf("Reacting to message event from channel %s", ev.Channel)
			message := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
			storeMsgError := storeMessage(message)
			if storeMsgError != nil {
				return "", storeMsgError
			}

			getSentimentError := getSentiment(message)
			if getSentimentError != nil {
				return "", getSentimentError
			}
		case *slackevents.AppMentionEvent:
			log.Printf("Reacting to app mention event from channel %s", ev.Channel)
			_, _, postError := postMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if postError != nil {
				message := fmt.Sprintf("Error while posting message %s", postError)
				log.Printf(message)
				return "", errors.New(message)
			}
		case *slackevents.AppHomeOpenedEvent:
			log.Println("Reacting to app home request event")
			homeViewForUser := slack.HomeTabViewRequest{
				Type:   "home",
				Blocks: apphome.UserHome("Cyril").Blocks,
			}
			homeViewAsJson, _ := json.Marshal(homeViewForUser)
			log.Printf("Sending view %s", homeViewAsJson)
			_, publishViewError := publishView(ev.User, homeViewForUser, ev.View.Hash)
			if publishViewError != nil {
				log.Println(publishViewError)
				return "", publishViewError
			}
		}
	}
	return challengeResponse, nil
}

// Slack Challenge is used to register the URL in the slack API config interface
// Should only be used once by slack when changing the events URL
func handleSlackChallenge(eventsAPIEvent slackevents.EventsAPIEvent, body []byte) (string, error) {
	var err error
	buf := new(bytes.Buffer)
	var r *slackevents.ChallengeResponse
	e := json.Unmarshal(body, &r)
	if e != nil {
		err = errors.New("Unable to register the URL")
		return "", err
	}
	buf.Write([]byte(r.Challenge))
	return buf.String(), err
}

var storeMessage = func(message *slackevents.MessageEvent) error {
	// Create DB
	tableName := os.Getenv("DYNAMODB_TABLE")
	dbError := dsedb.CreateDBIfNotCreated(tableName)
	if dbError {
		return errors.New("Database could not be created")
	}

	// Save in DB
	messageBytes, e := json.Marshal(message)
	if e != nil {
		return errors.New("Message could not be parsed before saving")
	}
	dbItem := dsedb.Message{
		UserId:         message.User,
		Text:           message.Text,
		CreatedAt:      message.EventTimeStamp.String(),
		SlackMessageId: message.EventTimeStamp.String(),
		SlackThreadId:  message.ThreadTimeStamp,
	}
	unmarshalError := json.Unmarshal(messageBytes, &dbItem)
	if unmarshalError != nil {
		return errors.New("Message could not JSONified")
	}
	log.Println(structs.Map(&dbItem))

	dbResult := dsedb.Store(tableName, structs.Map(&dbItem))
	if !dbResult {
		errorMsg := "Could not store message in DB"
		log.Println(errorMsg)
		return errors.New(errorMsg)
	}
	
	log.Println("Message was stored successfully")
	
	return nil
}

var getSentiment = func(message *slackevents.MessageEvent) error {
	tableName := os.Getenv("DYNAMODB_TABLE")
	apiKey := os.Getenv("PD_API_KEY")
	apiURL := os.Getenv("PD_API_URL")
	text := message.Text
	sentimentAnalysis, sentimentError := nlp.GetSentiment(text, apiURL, apiKey)
	if sentimentError != nil {
		errorMsg := "Could not analyze message"
		log.Println(errorMsg)
		return errors.New(errorMsg)
	}
	dbResult := dsedb.Update(tableName, message.EventTimeStamp.String(), sentimentAnalysis.Sentiment)
	if !dbResult {
		log.Println("Could not update message with sentiment")
	} else {
		log.Println("Message was updated successfully with sentiment")
	}
	return nil
}

