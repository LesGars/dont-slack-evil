package apphome

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/go-test/deep"
	"github.com/slack-go/slack"
)

func TestHomeBasicSections(t *testing.T) {
	expectedJson := heredoc.Doc(`
		{
			"type": "home",
			"blocks": [
				{
					"type": "section",
					"text": {
						"type": "mrkdwn",
						"text": "*Don't Slack Evil Performance*"
					},
					"accessory": {
						"type": "button",
						"text": {
							"type": "plain_text",
							"text": "Manage App Settings",
							"emoji": true
						},
						"value": "app_settings"
					}
				},
				{
					"type": "section",
					"text": {
						"type": "plain_text",
						"text": " :wave: Hello Cyril · find your DSE stats below",
						"emoji": true
					}
				},
				{
					"type": "divider"
				},
				{
					"type": "section",
					"fields": [
						{
							"type": "mrkdwn",
							"text": "*Current Quarter*\n(ends in 42 days)\nNumber of slack messages: 42\nEvil messages: 24\nImproved messages with DSE: 12/24"
						},
						{
							"type": "mrkdwn",
							"text": "*Top Channels with evil messages*\n:airplane: General · 30% (142)\n:taxi: Code Reviews · 66% (43)\n:knife_fork_plate: Direct Messages · 18% (75)"
						}
					]
				},
				{
					"type": "divider"
				}
			]
		}
	`)

	expectedBytes := []byte(expectedJson)
	var expectedObject slack.Message
	err := json.Unmarshal(expectedBytes, &expectedObject)
	if err != nil {
		fmt.Println("error:", err)
	}

	actual := slack.NewBlockMessage(HomeBasicSections("Cyril")...)
	actual.Msg.Type = "home"

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
