package main

import (
	"fmt"
	"encoding/json"
		// map[string]interface{}
		// []interface{}
	// "os"
	// "strings"
	// "math/rand"
	// "strconv"
	// "time"

	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"encoding/json"
)

type MessagesPendingApproval struct {
  MessagesPendingApproval []MessagePendingApproval `json:"messages"`
}

type MessagePendingApproval struct {
  Name string `json:"name"`
  EvilIndex string `json:"evilIndex"`
  Date int `json:"date"`
  MessageLink string `json:"messageLink"`
  Original string `json:"original"`
  DSESuggestion string `json:"dseSuggestion"`
}

func UserHome(userId string) map[string]interface{} {
	// TODO : get user name and Global stats from DynamoDB
	return map[string]interface{}{
		"type": "home",
		"blocks": HOME_HEADER_SECTIONS + helloSections(userId) + messageStats(userId)
	}
}

const HomeHeaderSections = [
	map[string]interface{}{
		"type": "section",
		"text": map[string]interface{}{
			"type": "mrkdwn",
			"text": "*Don't Slack Evil Performance*"
		},
		"accessory": map[string]interface{}{
			"type": "button",
			"text": map[string]interface{}{
				"type": "plain_text",
				"text": "Manage App Settings",
				"emoji": true
			},
			"value": "app_settings"
		}
	}
]

func helloSections(userId string) []interface{} {
	return [
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "plain_text",
				"text": fmt.Sprintf(" :wave: Hello %s · find your DSE stats below", translateUserIdToUserName(userId)),
				"emoji": true
			}
		},
		map[string]interface{}{
			"type": "divider"
		},
	]
}

func translateUserIdToUserName(userId string) string {
	// TODO
	return userId
}



func messageStats(userId string) []interface{} {
	return [
		map[string]interface{}{
			"type": "section",
			"fields": []interface{}{
				map[string]interface{}{
					"type": "mrkdwn",
					"text": "*Current Quarter*\n(ends in 53 days)\nNumber of slack messages: 42\nEvil messages: 24\nImproved messages with DSE: 12/24"
				},
				map[string]interface{}{
					"type": "mrkdwn",
					"text": "*Top Channels with evil messages*\n:airplane: General · 30% (142)\n:taxi: Code Reviews · 66% (43) \n:knife_fork_plate: Direct Messages · 18% (75)"
				}
			]
		},
		{
			"type": "context",
			"elements": []interface{}{
				map[string]interface{}{
					"type": "image",
					"image_url": "https://api.slack.com/img/blocks/bkb_template_images/placeholder.png",
					"alt_text": "placeholder"
				}
			]
		},
	]
}

func messagesPendingApprovalSections(messages string[]) []interface{}
	messagesFromJson := parseTestMessages()

	return [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Messages Awaiting Englightnment*"
			}
		},
		{
			"type": "divider"
		},
		] + mapTestJsonToMessage(messagesFromJson)
	}
}

func parseTestMessages() {
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var messages MessagesPendingApproval

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &messages)

	return messages

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	// for i := 0; i < len(users.Users); i++ {
	//     fmt.Println("User Type: " + users.Users[i].Type)
	//     fmt.Println("User Age: " + strconv.Itoa(users.Users[i].Age))
	//     fmt.Println("User Name: " + users.Users[i].Name)
	//     fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	// }
}

func mapTestJsonToMessage(messages string[]) map[string]interface{}{
	map[string]interface{}{
	{
		"type": "context",
		"elements": [
			{
				"type": "mrkdwn",
				"text": "Submitted by"
			},
			{
				"type": "image",
				"image_url": "https://api.slack.com/img/blocks/bkb_template_images/profile_3.png",
				"alt_text": "Evil Guy"
			},
			{
				"type": "mrkdwn",
				"text": "*Evil Guy*"
			}
		]
	},
	{
		"type": "section",
		"text": {
			"type": "mrkdwn",
			"text": "Channel : *#general*\nEstimated evil index: :imp::imp::imp: · *quite evil*\nDate: *2019/10/16 8:44*\n*<fakelink.toUrl.com|Link to message>*"
		},
		"accessory": {
			"type": "image",
			"image_url": "https://api.slack.com/img/blocks/bkb_template_images/creditcard.png",
			"alt_text": "credit card"
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
		"type": "actions",
		"elements": [
			{
				"type": "button",
				"text": {
					"type": "plain_text",
					"text": "Approve",
					"emoji": true
				},
				"style": "primary",
				"value": "approve"
			},
			{
				"type": "button",
				"text": {
					"type": "plain_text",
					"text": "Decline",
					"emoji": true
				},
				"style": "danger",
				"value": "decline"
			},
			{
				"type": "button",
				"text": {
					"type": "plain_text",
					"text": "View Message",
					"emoji": true
				},
				"value": "details"
			}
		]
	},
	{
		"type": "divider"
	}
}
