package apphome

import (
	"dont-slack-evil/helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/slack-go/slack"
)

type Message struct {
	UserName      string `json:"userName"`
	EvilIndex     string `json:"evilIndex"`
	Date          string `json:"date"`
	Channel       string `json:"channel"`
	MessageLink   string `json:"messageLink"`
	Original      string `json:"original"`
	DSESuggestion string `json:"dseSuggestion"`
}

func EnlightenmentSection() []slack.Block {
	messagesText := slack.NewTextBlockObject("mrkdwn", "*Messages Awaiting Englightnment*", false, false)
	messagesSection := slack.NewSectionBlock(messagesText, nil, nil)

	blocks := []slack.Block{messagesSection, slack.NewDividerBlock()}

	return append(blocks, EnlightenmentMessages()...)
}

func EnlightenmentMessages() []slack.Block {
	messages := parseTestMessages()
	var blocks []slack.Block

	for _, message := range messages {
		blocks = append(blocks, formatMessageToBeEnlightened(message)...)
	}
	return blocks
}

func parseTestMessages() []Message {
	file := path.Join(os.Getenv("GOPATH"), "apphome/sample.json")

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	log.Printf("Messages read from JSON: %s", byteValue)

	var messages []Message
	json.Unmarshal([]byte(byteValue), &messages)
	return messages
}

func formatMessageToBeEnlightened(message Message) []slack.Block {
	submittedBy := slack.NewTextBlockObject("plain_text", "Submitted by", false, false)
	submittedByImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_3.png", message.UserName)
	submittedByName := slack.NewTextBlockObject("plain_text", message.UserName, false, false)
	contextBlock := slack.NewContextBlock("", submittedBy, submittedByImage, submittedByName)

	estimatedEvilIndexText := slack.NewTextBlockObject("mrkdwn",
		heredoc.Doc(fmt.Sprintf(`
			Channel: %s
			Estimated evil index: %s
			Date: %s
			<%s|Link to message>`,
			message.Channel, message.EvilIndex, message.Date, message.MessageLink,
		)), false, false,
	)
	evilIndexBlock := slack.NewSectionBlock(estimatedEvilIndexText, nil, nil)

	dseSuggestionText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("Original\n%s\n\nKindly suggested alternative\n\n%s",
			helpers.QuoteForSlack(message.Original),
			helpers.QuoteForSlack(message.DSESuggestion),
		), false, false,
	)
	dseSuggestionBlock := slack.NewContextBlock("", dseSuggestionText)
	return []slack.Block{contextBlock, evilIndexBlock, dseSuggestionBlock, slack.NewDividerBlock()}
}
