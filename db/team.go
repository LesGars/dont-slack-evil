package db

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fatih/structs"
)

type Team struct {
	SlackTeamId            string `json:"slack_team_id"`
	SlackBotUserToken      string `json:"slack_bot_user_oauth_token"`
	SlackRegularOauthToken string `json:"slack_regular_oauth_token"`
}

func FindOrCreateTeamById(id string) (*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "-teams"
	// CreateTableIfNotCreated(tableName, "slack_team_id")

	team, err := FindTeamById(id)
	if err != nil {
		if err != nil {
			// TODO: intercept errorNotFound to create a team with empty bot token, and ask to complete OAuth flow
			if !Store(tableName, structs.New(Team{SlackTeamId: id}).Map()) {
				log.Printf("Error creating a team of ID %s", id)
			}
		} else {
			log.Printf("%s", err) // WIll print a not found OR other useful like access not granted
		}
		return nil, err
	}
	log.Printf("Found team of ID %s called %s", team.SlackTeamId, "TODO")
	return team, nil
}

func FindTeamById(id string) (*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "-teams"
	out, err := Get(tableName,
		map[string]*dynamodb.AttributeValue{
			"slack_team_id": {
				S: aws.String(id),
			},
		},
	)
	if err != nil {
		log.Printf("Error retrieving an item ID")
		return nil, err
	}

	var team Team
	unMarshallErr := dynamodbattribute.UnmarshalMap(out.Item, &team)
	if err != nil {
		log.Printf("Error unmarshalling a DynamoDB item ID as a Team")
		return nil, err
	}

	return &team, unMarshallErr
}
