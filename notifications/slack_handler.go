package notifications

import (
	"dont-slack-evil/apphome"
	"log"
	dsedb "dont-slack-evil/db"
	"github.com/slack-go/slack"
)

// SendNotifications loops through all the users and sends a message to those who sent too many
// messages of bad quality over the last quarter
func SendNotifications() (int, error) {
	notificationsSent := 0
	teams, teamsErr := dsedb.GetTeams()
	if teamsErr != nil {
		return 0, teamsErr
	}

	for _, team := range teams {
		slackBotUserApiClient := slack.New(team.SlackBotUserToken)
		users, err := slackBotUserApiClient.GetUsers()
		if err != nil {
			log.Printf("Could not instantiate bot client for team %v", team.SlackTeamId)
			continue
		}
		for _, user := range users {
			userId := user.ID
			tooManyBadQualityMessagesLastQuarter := apphome.HasTooManyBadQualityMessagesLastQuarter(userId)
			if tooManyBadQualityMessagesLastQuarter {
				conversationParameters := slack.OpenConversationParameters{
					Users: []string{userId},
				}
				channel, _, _, conversationErr := slackBotUserApiClient.OpenConversation(&conversationParameters)
				if conversationErr != nil {
					log.Printf("Could not open conversation for user %v: %v", userId, conversationErr)
				} else {
					slackBotUserApiClient.PostMessage(channel.ID, slack.MsgOptionText("Hello! It looks like you sent too many negative messages over the last quarter", false))
					notificationsSent++
				}
			}
		}
	}

	return notificationsSent, nil
}

// SendLeaderboardNotification sends the leaderboard notification
func SendLeaderboardNotification() (int, error) {
	log.Println("Hi Im here")
	return 0, nil
}
