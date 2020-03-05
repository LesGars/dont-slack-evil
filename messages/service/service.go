package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"dont-slack-evil/apphome"

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
var ParseEvent = slackevents.ParseEvent

// ParseEvent is the assignation of slackBotUserApiClient.PostMessage to a variable,
// in order to make it mockable
var PostMessage = slackBotUserApiClient.PostMessage

// PublishView is the assignation of slackBotUserApiClient.PublishView to a variable,
// in order to make it mockable
var PublishView = slackBotUserApiClient.PublishView 

// HandleEvent uses Slack's Event API to respond to an event emitted by our application
func HandleEvent(body []byte) (string, error) {
	var challengeResponse string
	eventsAPIEvent, e := ParseEvent(
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
		case *slackevents.AppMentionEvent:
			log.Printf("Reacting to app mention event from channel %s", ev.Channel)
			_, _, postError := PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
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
			_, publishViewError := PublishView(ev.User, homeViewForUser, ev.View.Hash)
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
