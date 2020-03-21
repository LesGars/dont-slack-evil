package apphome

import (
	dsedb "dont-slack-evil/db"
	"dont-slack-evil/leaderboard"
	"dont-slack-evil/stats"
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/slack-go/slack"
)

// Wraps the HomeSections with the top level slack payload
func UserHome(userId string, userName string, apiForTeam dsedb.ApiForTeam) slack.Message {
	message := slack.NewBlockMessage(
		HomeSections(userName, userId, apiForTeam)...,
	)
	message.Msg.Type = "home"

	return message
}

// Note this function
func HomeSections(userName string, userId string, apiForTeam dsedb.ApiForTeam) []slack.Block {
	return append(
		introSections(userName),
		statsSections(userId, apiForTeam)...,
	)
}

func introSections(userName string) []slack.Block {
	headerOptionsTxt := slack.NewTextBlockObject("plain_text", "Manage App Settings", true, false)
	headerButton := slack.NewButtonBlockElement("", "app_settings", headerOptionsTxt)
	headerText := slack.NewTextBlockObject("mrkdwn", "*Don't Slack Evil Performance*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, slack.NewAccessory(headerButton))

	helloText := slack.NewTextBlockObject("plain_text", fmt.Sprintf(" :wave: Hello %s · find your DSE stats below", userName), true, false)
	helloSection := slack.NewSectionBlock(helloText, nil, nil)

	return []slack.Block{headerSection, helloSection, slack.NewDividerBlock()}
}

func statsSections(userId string, apiForTeam dsedb.ApiForTeam) []slack.Block {
	stats := stats.HomeStatsForUser(userId)
	messageStats := slack.NewTextBlockObject("mrkdwn",
		heredoc.Doc(fmt.Sprintf(`
			*All time*
			Number of analyzed messages: %d
			Number of messages of good quality : %d

			Your overall positivity : %d%%

			*Current Quarter*
			(ends in %d days)
			Number of analyzed messages: %d
			Number of messages of good quality : %d

			Your positivity this quarter : %d%%`,
			stats.MessagesAnalyzedAllTime, stats.MessagesOfGoodQualityAllTime, int(stats.PercentageOfMessagesOfGoodQualityAllTime*100),
			42,
			stats.MessagesAnalyzedSinceQuarter, stats.MessagesOfGoodQualitySinceQuarter, int(stats.PercentageOfMessagesOfGoodQualitySinceQuarter*100),
		)),
		false, false,
	)

	leaderboards, err := leaderBoardsWithCurrentUser(apiForTeam, userId)
	var fields []*slack.TextBlockObject
	if err != nil {
		fields = []*slack.TextBlockObject{messageStats}
	} else {
		fields = []*slack.TextBlockObject{messageStats, leaderboards}
	}

	statsSection := slack.NewSectionBlock(nil, fields, nil)
	return []slack.Block{statsSection, slack.NewDividerBlock()}
}

func leaderBoardsWithCurrentUser(apiForTeam dsedb.ApiForTeam, userId string) (*slack.TextBlockObject, error) {
	scores, leaderboardErr := leaderboard.LeaderboardsForTeam(apiForTeam.SlackBotUserApiClient)
	if leaderboardErr != nil {
		return nil, leaderboardErr
	}
	if len(scores) < 3 {
		return nil, errors.New("Cannot add a leaderboard to the home: not enough people")
	}
	var isFirstPlace, isSecondPlace, isThirdPlace, isNthPlace string
	for i, _ := range scores {
		if scores[i].ID == userId {
			if i == 0 {
				isFirstPlace = " (you)"
			} else if i == 1 {
				isSecondPlace = " (you)"
			} else if i == 2 {
				isThirdPlace = " (you)"
			} else {
				isNthPlace = fmt.Sprintf(":face_with_monocle: You with a %.2f%% score (%d / %d)",
					scores[i].Score*100, scores[i].Good, scores[i].Total,
				)
			}
			break
		}
	}
	weeklyLeaderboardText := heredoc.Doc(fmt.Sprintf(`
		*Weekly positivity rankings:*

		Here are the standings for this quarter:
		:first_place_medal: <@%s>%s with a %.2f%% score (%d / %d)
		:second_place_medal: <@%s>%s with a %.2f%% score (%d / %d)
		:third_place_medal: <@%s>%s with a %.2f%% score (%d / %d)
		%s`, // would be painful to make this one conditional, it's ok to have a blank line with spaces
		scores[0].ID, isFirstPlace, scores[0].Score*100, scores[0].Good, scores[0].Total,
		scores[1].ID, isSecondPlace, scores[1].Score*100, scores[1].Good, scores[1].Total,
		scores[2].ID, isThirdPlace, scores[2].Score*100, scores[2].Good, scores[2].Total,
		isNthPlace,
	))

	return slack.NewTextBlockObject("mrkdwn", weeklyLeaderboardText, false, false), nil
}
