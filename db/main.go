package db

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func DynamoDBClient() *dynamodb.DynamoDB {
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{Profile: "dont-slack-evil-hackaton"},
			&credentials.EnvProvider{},
		},
	)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
	}))
	return dynamodb.New(sess)
}

// CreateTableIfNotCreated creates DynamoDB table if it doesn't exist
func CreateTableIfNotCreated(tableName string, mainKey string) error {
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(mainKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(mainKey),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}
	_, createTableErr := DynamoDBClient().CreateTable(createTableInput)
	if createTableErr != nil {
		if !strings.Contains(createTableErr.Error(), "Table already exists") {
			log.Println(createTableErr.Error())
			return createTableErr
		}
	} else {
		log.Println("Created the table", tableName)
	}
	return nil
}

// Store an item in the database
func Store(tableName string, item map[string]interface{}) bool {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("Got error marshalling new item:")
		log.Println(err.Error())
		return false
	}
	putItemInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = DynamoDBClient().PutItem(putItemInput)
	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())
		return false
	}
	return true
}

// Update an item in the database
func Update(tableName string, slackMessageId string, sentiment Sentiment) bool {
	expr, err := dynamodbattribute.MarshalMap(sentiment)
	if err != nil {
		log.Println("Got error marshalling info:")
		log.Println(err.Error())
		return false
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":updated_sentiment": {
				M: expr,
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"slack_message_id": {
				S: aws.String(slackMessageId),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set sentiment = :updated_sentiment"),
	}

	_, updateErr := DynamoDBClient().UpdateItem(input)
	if updateErr != nil {
		log.Println(updateErr.Error())
		return false
	}
	return true
}

// Get an item in the database
func Get(tableName string, key map[string]*dynamodb.AttributeValue) (*dynamodb.GetItemOutput, error) {
	return DynamoDBClient().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
}

// Convert a DynamoDB scan output to an integer
func ScanToInt(result *dynamodb.ScanOutput, err error) (int, error) {
	return int(*result.Count), nil
}

// Performe a Scan operation in DynamoDB
func Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	result, err := DynamoDBClient().Scan(input)
	if err != nil {
		log.Printf("Error during DynamoDB Scan")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return nil, err
	}
	return result, nil
}

func unmarshalScanResults(scanOutput *dynamodb.ScanOutput) []interface{} {
	var recs []interface{}

	resultsErr := dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &recs)
	if resultsErr != nil {
		log.Printf("failed to unmarshal Dynamodb Scan Items, %v", resultsErr)
	}
	return recs
}
