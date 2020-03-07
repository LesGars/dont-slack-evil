package apphome

import (
	dsedb "dont-slack-evil/db"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/jinzhu/now"
)

var thresholdQuality float64 = 0.2

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
	badQualityFilt := badQualityFilt()
	stats := DSEHomeStats{
		MessagesAnalyzedAllTime:     messagesAnalyzed(userIdFilt),
		MessagesOfBadQualityAllTime: messagesAnalyzed(expression.And(badQualityFilt, userIdFilt)),
	}
	if stats.MessagesAnalyzedAllTime != 0 {
		stats.PercentageOfMessagesOfBadQualityAllTime = float64(stats.MessagesOfBadQualityAllTime) / float64(stats.MessagesAnalyzedAllTime)
	}
	return stats
}

func userIdFilt(userId string) expression.ConditionBuilder {
	return expression.Equal(expression.Name("user_id"), expression.Value(userId))
}

func badQualityFilt() expression.ConditionBuilder {
	return expression.LessThan(expression.Name("sentiment.negative"), expression.Value(thresholdQuality))
}

func sinceBeginningOfQuarterFilt() expression.ConditionBuilder {
	return expression.GreaterThan(expression.Name("written_at"), expression.Value(now.BeginningOfQuarter()))
}

func messagesAnalyzed(userIdFilt expression.ConditionBuilder) int {
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
		TableName:                 aws.String(os.Getenv("DYNAMODB_TABLE_PREFIX") + "messages"),
	}

	val, err := dsedb.ScanToInt(dsedb.Scan(input))

	if err != nil {
		return 0
	}
	return val
}
