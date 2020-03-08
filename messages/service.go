package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"dont-slack-evil/apphome"
	dsedb "dont-slack-evil/db"
	"dont-slack-evil/nlp"

	"github.com/fatih/structs"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// If you need to send a message from the app's "bot user", use this bot client
var slackBotUserOauthToken = os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")
var slackBotUserApiClient = slack.New(slackBotUserOauthToken)

// ParseEvent is the assignation of slackevents.ParseEvent to a variable,
// in order to make it mockable
var parseEvent = slackevents.ParseEvent

// ParseEvent is the assignation of slackBotUserApiClient.PostMessage to a variable,
// in order to make it mockable
var postMessage = slackBotUserApiClient.PostMessage

// PublishView is the assignation of slackBotUserApiClient.PublishView to a variable,
// in order to make it mockable
var publishView = slackBotUserApiClient.PublishView

var getUserInfo = slackBotUserApiClient.GetUserInfo

var userHome = apphome.UserHome

func analyzeMessage(message *slackevents.MessageEvent, apiForTeam ApiForTeam) (string, error) {
	log.Printf("Reacting to message event from channel %s", message.Channel)
	storeMsgError := storeMessage(message, apiForTeam)
	if storeMsgError != nil {
		if !strings.Contains(storeMsgError.Error(), "Database could not be created") {
			return "", storeMsgError
		}
		log.Printf("Could not save initial message %s", storeMsgError)
	}

	getSentimentError := getSentiment(message)
	if getSentimentError != nil {
		return "", getSentimentError
	}

	return "", nil
}

func yesHello(message *slackevents.AppMentionEvent, apiForTeam ApiForTeam) (string, error) {
	log.Printf("Reacting to app mention event from channel %s", message.Channel)
	_, _, postError := apiForTeam.SlackBotUserApiClient.PostMessage(message.Channel, slack.MsgOptionText("Yes, hello.", false))
	if postError != nil {
		message := fmt.Sprintf("Error while posting message %s", postError)
		log.Printf(message)
		return "", errors.New(message)
	}
	return "", nil
}

func updateAppHome(ev *slackevents.AppHomeOpenedEvent, apiForTeam ApiForTeam) (string, error) {
	log.Println("Reacting to app home request event")
	userID := ev.User
	user, getUserInfoErr := getUserInfo(userID)
	if getUserInfoErr != nil {
		log.Println(getUserInfoErr)
	}
	userName := user.RealName
	homeViewForUser := slack.HomeTabViewRequest{
		Type:   "home",
		Blocks: userHome(userID, userName).Blocks,
	}
	homeViewAsJson, _ := json.Marshal(homeViewForUser)
	log.Printf("Sending view %s", homeViewAsJson)
	_, publishViewError := apiForTeam.SlackBotUserApiClient.PublishView(ev.User, homeViewForUser, ev.View.Hash)
	if publishViewError != nil {
		log.Println(publishViewError)
		return "", publishViewError
	}

	return "", nil
}

var storeMessage = func(message *slackevents.MessageEvent, apiForTeam ApiForTeam) error {
	// Create DB
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "messages"
	dbError := dsedb.CreateTableIfNotCreated(tableName, "slack_message_id")
	if dbError != nil {
		return dbError
	}

	// Save in DB
	messageBytes, e := json.Marshal(message)
	if e != nil {
		return errors.New("Message could not be parsed before saving")
	}
	timeUnix, _ := message.EventTimeStamp.Int64()
	log.Printf("Temps timestamp %s", message.EventTimeStamp)
	log.Printf("Temps unix %d", timeUnix)
	dbItem := dsedb.Message{
		UserId:         message.User,
		Text:           message.Text,
		CreatedAt:      time.Unix(timeUnix, 0).Format(time.RFC3339),
		SlackMessageId: message.EventTimeStamp.String(),
		SlackThreadId:  message.ThreadTimeStamp,
		SlackTeamId:    apiForTeam.Team.SlackTeamId,
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
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "messages"
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
