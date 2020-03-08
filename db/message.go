package db

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack/slackevents"
)

type Message struct {
	UserId         string    `json:"user_id"`
	SlackMessageId string    `json:"slack_message_id"`
	SlackThreadId  string    `json:"slack_thread_id"`
	SlackTeamId    string    `json:"slack_team_id"`
	Text           string    `json:"text"`
	Analyzed       bool      `json:"analyzed"`
	CreatedAt      string    `json:"created_at"`
	Quality        float64   `json:"quality"`
	Sentiment      Sentiment `json:"sentiment"`
}

type Sentiment struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}

func NewMessageFromSlack(message *slackevents.MessageEvent, teamId string) (*Message, error) {
	messageBytes, e := json.Marshal(message)
	if e != nil {
		return nil, errors.WithMessage(e, "Message could not be parsed before saving")
	}
	timeUnix, _ := message.EventTimeStamp.Float64()
	dbItem := Message{
		UserId:         message.User,
		Text:           message.Text,
		CreatedAt:      time.Unix(int64(timeUnix), 0).Format(time.RFC3339),
		SlackMessageId: message.EventTimeStamp.String(),
		SlackThreadId:  message.ThreadTimeStamp,
		SlackTeamId:    teamId,
	}
	unmarshalError := json.Unmarshal(messageBytes, &dbItem)
	if unmarshalError != nil {
		return nil, errors.WithMessage(unmarshalError, "Message could not JSONified")
	}
	return &dbItem, nil
}
