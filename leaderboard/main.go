package leaderboard

import (
	"errors"
	"fmt"
	"log"
	"sort"

	dsedb "dont-slack-evil/db"
	"dont-slack-evil/stats"

	"github.com/slack-go/slack"
)

// UserScore contains the positivity score of a Slack user
type UserScore struct {
	ID    string
	Name  string
	Good  int
	Total int
	Score float64
}

// SendLeaderboardNotification sends the leaderboard notification
func SendLeaderboardNotification() (int, error) {
	notificationsSent := 0
	teams, teamsErr := dsedb.GetTeams()
	log.Println("Analyzing weekly messages for", len(teams), "teams")
	if teamsErr != nil {
		log.Println(teamsErr)
		return 0, teamsErr
	}
	for _, team := range teams {
		slackBotUserApiClient := slack.New(team.SlackBotUserToken)
		log.Println("Analyzing weekly messages from team", team.SlackTeamId)
		userScores, leaderboardErr := LeaderboardsForTeam(slackBotUserApiClient)
		if leaderboardErr != nil {
			continue
		}
		if sendWebhookNotification(team, userScores) == nil {
			notificationsSent++
		}
	}

	// Note : here we decide not to throw an error if some of the teams had errors (as they may have outdated tokens, etc)
	// But we'll have to watch this
	if notificationsSent == 0 {
		return 0, errors.New("No message was sent, this most likely indicates a problem 🧐")
	} else {
		return notificationsSent, nil
	}
}

func LeaderboardsForTeam(slackBotUserApiClient dsedb.SlackApiInterface) ([]UserScore, error) {
	users, err := slackBotUserApiClient.GetUsers()
	var userScores []UserScore
	if err != nil {
		log.Printf("Could not instantiate bot client for team %v", slackBotUserApiClient)
		return userScores, err
	}
	log.Println("Found", len(users), "users in team")
	for _, user := range users {
		// This is the best way I found to distinguish bots from real users
		// Note that user.IsBot doesn't work because it's false even for some bot users...
		if len(user.Profile.BotID) != 0 || user.Name == "slackbot" || user.IsBot {
			continue
		}

		good, total := stats.GetWeeklyStats(user.ID)
		var score float64
		if total > 0 {
			score = float64(good) / float64(total)
		} else {
			score = 0
		}
		userScore := UserScore{
			ID:    user.ID,
			Name:  user.RealName,
			Good:  good,
			Total: total,
			Score: score,
		}
		userScores = append(userScores, userScore)
	}

	// Sort by positivity scores
	sort.Slice(userScores, func(i, j int) bool {
		return userScores[i].Score > userScores[j].Score
	})

	return userScores, nil
}

func sendWebhookNotification(team *dsedb.Team, userScores []UserScore) error {
	var text = "*Weekly positivity rankings:*"
	if len(userScores) > 0 {
		text += fmt.Sprintf(
			"\n\nCongratulations to <@%s> for being the most positive person this week :tada:",
			userScores[0].ID,
		)
		text += "\n\nHere are the standings:"
		text += fmt.Sprintf(
			"\n:first_place_medal: <@%s> with a %.2f score (%d / %d)",
			userScores[0].ID,
			userScores[0].Score*100,
			userScores[0].Good,
			userScores[0].Total,
		)
	}
	if len(userScores) > 1 {
		text += fmt.Sprintf(
			"\n:second_place_medal: <@%s> with a %.2f score (%d / %d)",
			userScores[1].ID,
			userScores[1].Score*100,
			userScores[1].Good,
			userScores[1].Total,
		)
	}
	if len(userScores) > 2 {
		text += fmt.Sprintf(
			"\n:third_place_medal: <@%s> with a %.2f score (%d / %d)",
			userScores[2].ID,
			userScores[2].Score*100,
			userScores[2].Good,
			userScores[2].Total,
		)
	}
	msg := slack.WebhookMessage{
		Channel: team.IncomingWebhook.ChannelID,
		Text:    text,
	}
	webhookError := slack.PostWebhook(team.IncomingWebhook.URL, &msg)
	if webhookError != nil {
		log.Println(webhookError)
		return webhookError
	} else {
		log.Println("Sending message to channel", team.IncomingWebhook.Channel, "for team", team.SlackTeamId)
	}
	return nil
}
