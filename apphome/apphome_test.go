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
							"text": "*All time*\nNumber of analyzed messages: 0\nNumber of messages of bad quality : 0\n% of messages of bad quality : 0.000000\n*Current Quarter*\n(ends in 42 days)\nNumber of analyzed messages: 0\nNumber of messages of bad quality : 0\n% of messages of bad quality : 0.000000"
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

	actual := slack.NewBlockMessage(HomeBasicSections("Cyril", "42")...)
	actual.Msg.Type = "home"

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
