package notifications

import (
	"dont-slack-evil/apphome"
	"log"
	"os"

	"github.com/slack-go/slack"
)

var slackBotUserOauthToken = os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")
var slackBotUserApiClient = slack.New(slackBotUserOauthToken)

func SendNotifications() error {

	userIds, err := apphome.UsersWithBadQualityMessages()
	if err != nil {
		return err
	}
	for _, userId := range userIds {
		conversationParameters := slack.OpenConversationParameters{
			Users: []string{userId},
		}
		channel, _, _, conversationErr := slackBotUserApiClient.OpenConversation(&conversationParameters)
		if conversationErr != nil {
			log.Printf("Could not open conversation for user %v: %v", userId, conversationErr)
		} else {
			slackBotUserApiClient.PostMessage(channel.ID, slack.MsgOptionText("Test", false))
		}
		
	}

	return nil
}