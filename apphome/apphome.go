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
	helloText := slack.NewTextBlockObject("plain_text", fmt.Sprintf(" :wave: Hello %s · find your DSE stats below", translateUserIdToUserName(userId)), true, false)
	hellloSection := slack.NewSectionBlock(helloText, nil, nil)

	// Build Message with blocks created above
	message := slack.NewBlockMessage(
		headerSection,
		hellloSection,
		divSection,
	)
	message.Msg.Type = "home"

	return message
}

func translateUserIdToUserName(userId string) string {
	// TODO in https://github.com/gjgd/dont-slack-evil/issues/16
	return userId
}
