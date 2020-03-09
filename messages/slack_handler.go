package messages

import (
	"bytes"
	"dont-slack-evil/db"
	dsedb "dont-slack-evil/db"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")

type ApiForTeam struct {
	Team                  db.Team
	SlackBotUserApiClient SlackApiInterface
}

type SlackApiInterface interface {
	// This interface is meant to make a *slack.Client mockable easily
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	PublishView(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error)
}

// SlackHandler uses Slack's Event API to respond to an event emitted by our application
func SlackHandler(body []byte) (string, error) {
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
		return handleSlackChallenge(eventsAPIEvent, body)
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent || eventsAPIEvent.Type == slackevents.AppMention {
		// Retrieve team data (token, etc)
		apiForTeam, apiForTeamErr := buildApiForTeam(eventsAPIEvent)
		if apiForTeamErr != nil {
			log.Printf("Could not retrieve team data: %s", apiForTeamErr)
			return "", apiForTeamErr
		}
		handleSlackEvent(eventsAPIEvent, *apiForTeam)
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

func handleSlackEvent(eventsAPIEvent slackevents.EventsAPIEvent, apiForTeam ApiForTeam) (string, error) {
	innerEvent := eventsAPIEvent.InnerEvent

	// Process the event using team data
	log.Printf("Processing an event of inner data %s", innerEvent.Data)
	switch ev := innerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		return analyzeMessage(ev, apiForTeam)
	case *slackevents.AppMentionEvent:
		return yesHello(ev, apiForTeam)
	case *slackevents.AppHomeOpenedEvent:
		return updateAppHome(ev, apiForTeam)
	}
	return "", nil
}

var buildApiForTeam = func(eventsAPIEvent slackevents.EventsAPIEvent) (*ApiForTeam, error) {
	// Retrieve team data (token, etc)
	team, teamErr := dsedb.FindOrCreateTeamById(eventsAPIEvent.TeamID)
	if teamErr != nil {
		log.Printf("Could not retrieve team data: %s", teamErr)
		return nil, teamErr
	}
	slackBotUserApiClient := slack.New(team.SlackBotUserToken)
	return &ApiForTeam{Team: *team, SlackBotUserApiClient: slackBotUserApiClient}, nil
}
