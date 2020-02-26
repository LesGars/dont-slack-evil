package apphome

import (
	"fmt"

	"github.com/slack-go/slack"
)

func UserHome(userId string) slack.Message {
	// Shared Assets for example
	divSection := slack.NewDividerBlock()

	// Header Section
	headerOptionsTxt := slack.NewTextBlockObject("plain_text", "Manage App Settings", true, false)
	headerButton := slack.NewButtonBlockElement("", "app_settings", headerOptionsTxt)
	headerText := slack.NewTextBlockObject("mrkdwn", "*Don't Slack Evil Performance*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, slack.NewAccessory(headerButton))

	// Hello
	helloText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(" :wave: Hello %s · find your DSE stats below", translateUserIdToUserName(userId)), false, false)
	hellloSection := slack.NewSectionBlock(helloText, nil, nil)

	// Build Message with blocks created above
	return slack.NewBlockMessage(
		headerSection,
		hellloSection,
		divSection,
	)
}

func translateUserIdToUserName(userId string) string {
	// TODO
	return userId
}
