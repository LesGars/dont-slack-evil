package apphome

import (
	"dont-slack-evil/db"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-test/deep"
	"github.com/slack-go/slack"
)

func TestHomeBasicSections(t *testing.T) {
	expectedMedals := strings.Replace(heredoc.Docf(`
		*Weekly positivity rankings:*

		Here are the standings for this quarter:
		:first_place_medal: <@UU7KH0J0P> with a 39%% score
		:second_place_medal: <@UTU9SCT6X> with a 38%% score
		:third_place_medal: <@UTT0779FC> with a 36%% score`,
	), "\n", `\n`, -1)
	expectedScores := strings.Replace(heredoc.Docf(`
		*All time*
		Number of analyzed messages: 0
		Number of messages of good quality : 0

		Your overall positivity : 0%%

		*Current Quarter*
		(ends in 42 days)
		Number of analyzed messages: 0
		Number of messages of good quality : 0

		Your positivity this quarter : 0%%`,
	), "\n", `\n`, -1)

	expectedJson := heredoc.Docf(`
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
							"text": "%s"
						},
						{
							"type": "mrkdwn",
							"text": "%s"
						}
					]
				},
				{
					"type": "divider"
				}
			]
		}
	`, expectedScores, expectedMedals)

	expectedBytes := []byte(expectedJson)
	var expectedObject slack.Message
	err := json.Unmarshal(expectedBytes, &expectedObject)
	if err != nil {
		fmt.Println("error:", err)
	}

	actual := slack.NewBlockMessage(HomeSections("Cyril", "42", db.ApiForTeam{})...)
	actual.Msg.Type = "home"

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
