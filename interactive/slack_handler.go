package interactive

import (
	"encoding/json"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func SlackHandler(payload string, response Response) error {
	var interaction slack.InteractionCallback

	err := json.Unmarshal([]byte(payload), &interaction)
	if err != nil {
		log.Printf("Could not parse action response JSON: %v", err)
		return err
	}
	// See https://github.com/slack-go/slack/blob/510942f19cfde364380379e47f5b4f780b5fb477/interactions.go
	switch interaction.Type {
	case slack.InteractionTypeMessageAction:
		return acknowledge(response)
	}

	return nil
}

func acknowledge(response Response) error {
	ack := slackevents.MessageActionResponse{
		ResponseType:    slack.ResponseTypeEphemeral,
		ReplaceOriginal: false,
		Text:            "",
	}
	ackMarshalled, err := json.Marshal(ack)
	if err != nil {
		return err
	}
	response.StatusCode = 200
	response.Body = string(ackMarshalled)
	return nil
}
