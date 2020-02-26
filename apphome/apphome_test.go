package apphome

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/go-test/deep"
	"github.com/slack-go/slack"
)

func TestUserHome(t *testing.T) {
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

	var actual slack.Message = UserHome("Cyril")

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
