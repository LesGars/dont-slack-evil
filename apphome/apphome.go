package apphome

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/slack-go/slack"
)

func UserHome(userId string, userName string) slack.Message {
	message := slack.NewBlockMessage(
		append(
			HomeBasicSections(userName, userId),
		)...,
	)
	message.Msg.Type = "home"

	return message
}

func HomeBasicSections(userName string, userId string) []slack.Block {
	return append(
		introSections(userName),
		statsSections(userId)...,
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

func statsSections(userId string) []slack.Block {
	stats := HomeStatsForUser(userId)
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

	weeklyLeaderboardText := heredoc.Doc(fmt.Sprintf(`
		*Weekly positivity rankings:*

		Here are the standings for this quarter:
		:first_place_medal: <@UU7KH0J0P> with a %1.f%% score
		:second_place_medal: <@UTU9SCT6X> with a %1.f%% score
		:third_place_medal: <@UTT0779FC> with a %1.f%% score`,
		0.391304347826087*100, 0.375*100, 0.356789*100))
	weeklyLeaderboard := slack.NewTextBlockObject("mrkdwn", weeklyLeaderboardText, false, false)

	fields := []*slack.TextBlockObject{messageStats, weeklyLeaderboard}
	statsSection := slack.NewSectionBlock(nil, fields, nil)
	return []slack.Block{statsSection, slack.NewDividerBlock()}
}
