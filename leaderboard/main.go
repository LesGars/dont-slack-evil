package leaderboard

import (
	"fmt"
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
			Good int
			Total int
			Score float64
		}
		var userScores []UserScore;
		for _, user := range users {
			// This is the best way I found to distinguish bots from real users
			// Note that user.IsBot doesn't work because it's false even for bot users...
			log.Println(user.RealName)
			log.Println(user.Profile.BotID)
			if (len(user.Profile.BotID) == 0) {
				good, total := apphome.GetWeeklyStats(user.ID)
				var score float64;
				if (total > 0) {
					score = float64(good) / float64(total)
				} else {
					score = 0
				}
				userScore := UserScore{
					ID: user.ID,
					Name: user.RealName,
					Good: good,
					Total: total,
					Score: score,
				}
				userScores = append(userScores, userScore)
			}
		}
		// Sort by positivity scores
		sort.Slice(userScores, func(i, j int) bool {
			return userScores[i].Score > userScores[j].Score
		})
		log.Println(userScores)
		top3 := userScores[:3]
		log.Println(top3)
	}

	return notificationsSent, nil
}