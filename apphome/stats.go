package apphome

import (
	dsedb "dont-slack-evil/db"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var thresholdQuality float64 = 0.5

type DSEHomeStats struct {
	MessagesAnalyzedAllTime                 int     `json:"messagesAnalyzedAllTime"`
	MessagesOfBadQualityAllTime             int     `json:"messagesOfBadQualityAllTime"`
	PercentageOfMessagesOfBadQualityAllTime float64 `json:"percentageOfMessagesOfBadQualityAllTime"`

	MessagesAnalyzedLastQuarter                 int     `json:"messagesAnalyzedLastQuarter"`
	MessagesOfBadQualityLastQuarter             int     `json:"messagesOfBadQualityLastQuarter"`
	PercentageOfMessagesOfBadQualityLastQuarter float64 `json:"percentageOfMessagesOfBadQualityLastQuarter"`
}

func HomeStatsForUser(userId string) DSEHomeStats {
	userIdFilt := userIdFilt(userId)
	return DSEHomeStats{
		MessagesAnalyzedAllTime: messagesAnalyzed(userIdFilt),
	}
}

func userIdFilt(userId string) expression.ConditionBuilder {
	return expression.Equal(expression.Name("user_id"), expression.Value(userId))
}

func messagesAnalyzed(userIdFilt expression.ConditionBuilder) int {
	// filt := expression.And(
	// 	// expression.GreaterThanEqual(expression.Name("quality"), expression.Value(thresholdQuality)),
	// 	userIdFilt,
	// )
	expr, buildErr := expression.NewBuilder().WithFilter(userIdFilt).Build()
	if buildErr != nil {
		log.Println("Got error building expression:")
		log.Println(buildErr.Error())
		return 42
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(os.Getenv("DYNAMODB_TABLE")),
	}

	val, err := dsedb.ScanToInt(dsedb.Scan(input))

	if err != nil {
		return 0
	}
	return val
}
