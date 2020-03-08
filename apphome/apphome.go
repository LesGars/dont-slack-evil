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
			EnlightenmentSection()...,
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
			Number of messages of bad quality : %d
			%% of messages of bad quality : %f
			*Current Quarter*
			(ends in %d days)
			Number of analyzed messages: %d
			Number of messages of bad quality : %d
			%% of messages of bad quality : %f`,
			stats.MessagesAnalyzedAllTime, stats.MessagesOfBadQualityAllTime, stats.PercentageOfMessagesOfBadQualityAllTime,
			42,
			stats.MessagesAnalyzedSinceQuarter, stats.MessagesOfBadQualitySinceQuarter, stats.PercentageOfMessagesOfBadQualitySinceQuarter,
		)),
		false, false,
	)

	fields := []*slack.TextBlockObject{messageStats}
	statsSection := slack.NewSectionBlock(nil, fields, nil)
	return []slack.Block{statsSection, slack.NewDividerBlock()}
}
