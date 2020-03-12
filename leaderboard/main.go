package leaderboard

import (
	"log"
	"sort"

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
		type UserScore struct {
			ID string
			Name string
			Score float64
		}
		var userScores []UserScore;
		for _, user := range users {
			// This is the best way I found to distinguish bots from real users
			// Note that user.IsBot doesn't work because it's false even for bot users...
			if (len(user.Profile.BotID) == 0) {
				userScore := UserScore{
					ID: user.ID,
					Name: user.RealName,
					Score: apphome.GetWeeklyPositivityScore(user.ID),
				}
				userScores = append(userScores, userScore)
			}
		}
		sort.Slice(userScores, func(i, j int) bool {
			return userScores[i].Score > userScores[j].Score
		})
		log.Println(userScores)
	}

	return notificationsSent, nil
}