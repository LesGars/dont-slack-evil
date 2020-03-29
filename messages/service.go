package messages

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"

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

var userHome = apphome.UserHome

var warnIfTooBadThreshold float64 = 0.33 // If above this value, send a message to the user

func analyzeMessageAndWarnIfTooNegative(message *slackevents.MessageEvent, apiForTeam dsedb.ApiForTeam) (string, error) {
	sentiment, err := analyzeMessage(message, apiForTeam)
	if err != nil {
		return "", err
	}

	return "", warnIfTooNegative(message, apiForTeam, *sentiment)
}

func analyzeMessage(message *slackevents.MessageEvent, apiForTeam dsedb.ApiForTeam) (*dsedb.Sentiment, error) {
	log.Printf("Reacting to message event from channel %s", message.Channel)
	storeMsgError := storeMessage(message, apiForTeam)
	if storeMsgError != nil {
		if !strings.Contains(storeMsgError.Error(), "Database could not be created") {
			return nil, storeMsgError
		}
		log.Printf("Could not save initial message %s", storeMsgError)
	}

	sentiment, getSentimentError := getSentiment(message)
	if getSentimentError != nil {
		return nil, getSentimentError
	}

	return sentiment, nil
}

func yesHello(message *slackevents.AppMentionEvent, apiForTeam dsedb.ApiForTeam) (string, error) {
	log.Printf("Reacting to app mention event from channel %s", message.Channel)
	_, _, postError := apiForTeam.SlackBotUserApiClient.PostMessage(message.Channel, slack.MsgOptionText("Yes, hello.", false))
	if postError != nil {
		message := fmt.Sprintf("Error while posting message %s", postError)
		log.Printf(message)
		return "", errors.New(message)
	}
	return "", nil
}

func updateAppHome(ev *slackevents.AppHomeOpenedEvent, apiForTeam dsedb.ApiForTeam) (string, error) {
	log.Println("Reacting to app home request event")
	userID := ev.User
	var userName string
	user, getUserInfoErr := apiForTeam.SlackBotUserApiClient.GetUserInfo(userID)
	if getUserInfoErr != nil {
		log.Printf("Error getting user ID %+v", getUserInfoErr)
		// Fallback to empty string
	} else {
		userName = user.RealName
	}

	homeViewForUser := slack.HomeTabViewRequest{
		Type:   "home",
		Blocks: userHome(userID, userName, apiForTeam).Blocks,
	}
	homeViewAsJson, _ := json.Marshal(homeViewForUser)
	log.Printf("Sending view %s", homeViewAsJson)
	_, publishViewError := apiForTeam.SlackBotUserApiClient.PublishView(ev.User, homeViewForUser, ev.View.Hash)
	if publishViewError != nil {
		log.Printf("Error updating the app home: %+v", publishViewError)
		return "", publishViewError
	}

	return "", nil
}

var storeMessage = func(message *slackevents.MessageEvent, apiForTeam dsedb.ApiForTeam) error {
	// Create DB
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "messages"
	dbError := dsedb.CreateTableIfNotCreated(tableName, "slack_message_id")
	if dbError != nil {
		return dbError
	}

	// Save in DB
	dbItem, dbItemErr := dsedb.NewMessageFromSlack(message, apiForTeam.Team.SlackTeamId)
	if dbItemErr != nil {
		return errors.WithMessage(dbItemErr, "Could not instanciate a new Message form slack")
	}

	dbResult := dsedb.Store(tableName, structs.Map(&dbItem))
	if !dbResult {
		errorMsg := "Could not store message in DB"
		log.Println(errorMsg)
		return errors.New(errorMsg)
	}

	log.Println("Message was stored successfully")

	return nil
}

var getSentiment = func(message *slackevents.MessageEvent) (*dsedb.Sentiment, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "messages"
	apiKey := os.Getenv("PD_API_KEY")
	apiURL := os.Getenv("PD_API_URL")
	text := message.Text
	sentimentAnalysis, sentimentError := nlp.GetSentiment(text, apiURL, apiKey)
	if sentimentError != nil {
		errorMsg := "Could not analyze message"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}
	dbResult := dsedb.Update(tableName, message.EventTimeStamp.String(), sentimentAnalysis.Sentiment)
	if !dbResult {
		log.Println("Could not update message with sentiment")
	} else {
		log.Println("Message was updated successfully with sentiment")
	}
	return &sentimentAnalysis.Sentiment, nil
}

var warnIfTooNegative = func(message *slackevents.MessageEvent, apiForTeam dsedb.ApiForTeam, sentiment dsedb.Sentiment) error {
	log.Printf("Warning of a bad message, detected negativity is %f", sentiment.Negative)
	if sentiment.Negative >= warnIfTooBadThreshold {
		messageVisibleOnlyByUser := heredoc.Docf(`
			:warning: Be careful ! Your message is %.0f%% negative

			:sunny: Try to stay positive to boost productivity and friendliness in the workspace

			:slack: Come talk to me to find out more!`, sentiment.Negative*100,
		)

		// This is annoying : an Ephemeral message is NOT notified to the user
		// So if it is the first message sent under a thread, the user would never see it
		// (Although is is visible when clicking the "start a thread" button in Slack)
		// Therefore, only send the message in a thread if it is not the first message to be threaded
		ephemeralMsgOptions := []slack.MsgOption{
			slack.MsgOptionText(messageVisibleOnlyByUser, false),
		}
		if message.ThreadTimeStamp != "" {
			ephemeralMsgOptions = append(ephemeralMsgOptions, slack.MsgOptionTS(message.ThreadTimeStamp))
		}
		_, postError := apiForTeam.SlackBotUserApiClient.PostEphemeral(
			message.Channel,
			message.User,
			ephemeralMsgOptions...,
		)
		log.Println("The message is too negative, sending a warning to the user")
		if postError != nil {
			message := fmt.Sprintf("Error while posting ephemeral message %s", postError)
			return errors.New(message)
		}
	}
	return nil
}
