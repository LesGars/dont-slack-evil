package interactive

import (
	"encoding/json"
	"log"

	"github.com/slack-go/slack"
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
		acknowledge(response)
		return nil
	}

	return nil
}

func acknowledge(response Response) {
	response.StatusCode = 200
	response.Body = "ok"
	response.Headers["Content-Type"] = "text"
}
