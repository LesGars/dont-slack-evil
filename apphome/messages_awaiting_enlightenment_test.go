package apphome

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/go-test/deep"
	"github.com/slack-go/slack"
)

func TestEnlightenmentMessages(t *testing.T) {
	expectedJson := heredoc.Doc(`
		{
			"type": "home",
			"blocks": [
				{
					"type": "context",
					"elements": [
						{
							"type": "plain_text",
							"text": "Submitted by"
						},
						{
							"type": "image",
							"image_url": "https://api.slack.com/img/blocks/bkb_template_images/profile_3.png",
							"alt_text": "Evil Guy"
						},
						{
							"type": "plain_text",
							"text": "Evil Guy"
						}
					]
				},
				{
					"type": "section",
					"text": {
						"type": "mrkdwn",
						"text": "Channel: #general\nEstimated evil index: :imp::imp::imp: · *quite evil*\nDate: 2019/10/16 8:44\n<https://lesgarshack.slack.com/archives/CTU9UVDKM/p1582569122004300|Link to message>"
					}
				},
				{
					"type": "context",
					"elements": [
						{
							"type": "mrkdwn",
							"text": "Original\n> Hey, are you really kidding me ??! Seriously this code is so shitty even my 5-year-old son could have done better !\n\nKindly suggested alternative\n\n> Hey, I'm a bit surprised by this :smiley:. This code feels like it could be improved without too much effort."
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

	actual := slack.NewBlockMessage(EnlightenmentMessages()...)
	actual.Msg.Type = "home"

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
