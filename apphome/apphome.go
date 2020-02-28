package apphome

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/slack-go/slack"
)

func UserHome(userId string) slack.Message {
	name := translateUserIdToUserName(userId)
	message := slack.NewBlockMessage(
		append(
			HomeBasicSections(name),
			EnlightenmentSection()...,
		)...,
	)
	message.Msg.Type = "home"

	return message
}

func HomeBasicSections(userName string) []slack.Block {
	return append(
		introSections(userName),
		statsSections()...,
	)
}

func translateUserIdToUserName(userId string) string {
	// TODO in https://github.com/gjgd/dont-slack-evil/issues/16
	return userId
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

func statsSections() []slack.Block {
	numberOfEvilMessages := 24
	numberOfImprovedMessages := 12
	numberOfSlackMessages := 42
	daysLeftUntilQuarterEnd := 42
	messageStats := slack.NewTextBlockObject("mrkdwn",
		heredoc.Doc(fmt.Sprintf(`
			*Current Quarter*
			(ends in %d days)
			Number of slack messages: %d
			Evil messages: %d
			Improved messages with DSE: %d/%d`,
			daysLeftUntilQuarterEnd, numberOfSlackMessages, numberOfEvilMessages, numberOfImprovedMessages, numberOfEvilMessages,
		)),
		false, false,
	)

	topChannelsText := slack.NewTextBlockObject("mrkdwn",
		heredoc.Doc(fmt.Sprintf(`
			*Top Channels with evil messages*
			:airplane: General · 30%% (142)
			:taxi: Code Reviews · 66%% (43)
			:knife_fork_plate: Direct Messages · 18%% (75)`,
		)),
		false, false,
	)
	fields := []*slack.TextBlockObject{messageStats, topChannelsText}
	statsSection := slack.NewSectionBlock(nil, fields, nil)
	return []slack.Block{statsSection, slack.NewDividerBlock()}
}
