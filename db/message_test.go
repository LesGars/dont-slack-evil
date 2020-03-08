package db

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/slack-go/slack/slackevents"
)

// If parseEvent fails, the handler should return an error
func TestNewMessageFromSlack(t *testing.T) {
	actual, _ := NewMessageFromSlack(&slackevents.MessageEvent{
		User:           "42",
		Text:           "blabla",
		TimeStamp:      "1583708649.001100",
		EventTimeStamp: json.Number("1583708649.001100"),
	}, "LesGarsHack")

	expected := Message{
		UserId:         "42",
		SlackMessageId: "1583708649.001100",
		SlackThreadId:  "",
		SlackTeamId:    "LesGarsHack",
		Text:           "blabla",
		Analyzed:       false,
		CreatedAt:      "2020-03-09T00:04:09+01:00",
		Quality:        0,
		Sentiment:      Sentiment{},
	}
	if diff := deep.Equal(expected, *actual); diff != nil {
		t.Error(diff)
	}
}
