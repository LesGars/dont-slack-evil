package leaderboard

import (
	"log"

	"dont-slack-evil/apphome"
	dsedb "dont-slack-evil/db"
	"github.com/slack-go/slack"
)

// SendLeaderboardNotification sends the leaderboard notification
func SendLeaderboardNotification() (int, error) {
	notificationsSent := 0
	teams, teamsErr := dsedb.GetTeams()
	if teamsErr != nil {
		log.Println(teamsErr)
		return 0, teamsErr
	}
	for _, team := range teams {
		// FIXME: remove
		if (team.SlackTeamId != "TU7KB9FB9") {
			continue
		}
		slackBotUserApiClient := slack.New(team.SlackBotUserToken)
		users, err := slackBotUserApiClient.GetUsers()
		if err != nil {
			log.Printf("Could not instantiate bot client for team %v", team.SlackTeamId)
			continue
		}
		log.Println(team.SlackTeamId)
		type UserScore struct {
			userID string
			userScore float64
		}
		var userScores []UserScore;
		for _, user := range users {
			// This is the best way I found to distinguish bots from real users
			// Note that user.IsBot doesn't work because it's false even for bot users...
			if (len(user.Profile.BotID) == 0) {
				log.Println(user.RealName)
				log.Println(user.ID)
			}
		}
	}

	return notificationsSent, nil
}