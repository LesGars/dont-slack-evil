package db

import (
	"errors"
	"log"
	"os"
	"strings"

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
	// CreateTableIfNotCreated(tableName, "slack_team_id")
	team, findErr := FindTeamById(id)
	if findErr != nil {
		// TODO: check the error string, I wasn't able to make sure of this one
		if !strings.Contains(findErr.Error(), "Item does not exist") {
			return createTeamById(id)
		} else {
			log.Printf("%s", findErr)
			return nil, findErr
		}
	}
	log.Printf("Found team of ID %s called %s", team.SlackTeamId, "TODO")
	return team, nil
}

func createTeamById(id string) (*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "teams"
	team := &Team{SlackTeamId: id}
	if Store(tableName, structs.New(team).Map()) {
		return team, nil
	} else {
		log.Printf("Error creating a team of ID %s", id)
		return nil, errors.New("Could not create the team in DynamoDB")
	}
}

func FindTeamById(id string) (*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "teams"
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
