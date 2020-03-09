package apphome

import (
	dsedb "dont-slack-evil/db"
	"log"
	"os"
	"time"

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

	MessagesAnalyzedSinceQuarter                 int     `json:"messagesAnalyzedLastQuarter"`
	MessagesOfBadQualitySinceQuarter             int     `json:"messagesOfBadQualityLastQuarter"`
	PercentageOfMessagesOfBadQualitySinceQuarter float64 `json:"percentageOfMessagesOfBadQualityLastQuarter"`
}

func HomeStatsForUser(userId string) DSEHomeStats {
	userIdFilt := userIdFilt(userId)
	badQualityFilt := badQualityFilt()
	sinceBeginningOfQuarterFilt := sinceBeginningOfQuarterFilt()
	stats := DSEHomeStats{
		MessagesAnalyzedAllTime:     messagesAnalyzed(userIdFilt),
		MessagesOfBadQualityAllTime: messagesAnalyzed(expression.And(badQualityFilt, userIdFilt)),

		MessagesAnalyzedSinceQuarter:     messagesAnalyzed(expression.And(userIdFilt, sinceBeginningOfQuarterFilt)),
		MessagesOfBadQualitySinceQuarter: messagesAnalyzed(expression.And(badQualityFilt, userIdFilt, sinceBeginningOfQuarterFilt)),
	}
	if stats.MessagesAnalyzedAllTime != 0 {
		stats.PercentageOfMessagesOfBadQualityAllTime = float64(stats.MessagesOfBadQualityAllTime) / float64(stats.MessagesAnalyzedAllTime)
	}
	if stats.MessagesAnalyzedSinceQuarter != 0 {
		stats.PercentageOfMessagesOfBadQualitySinceQuarter = float64(stats.MessagesOfBadQualitySinceQuarter) / float64(stats.MessagesAnalyzedSinceQuarter)
	}
	return stats
}

func userIdFilt(userId string) expression.ConditionBuilder {
	return expression.Equal(expression.Name("user_id"), expression.Value(userId))
}

func badQualityFilt() expression.ConditionBuilder {
	return expression.GreaterThan(expression.Name("sentiment.negative"), expression.Value(thresholdQuality))
}

func sinceBeginningOfQuarterFilt() expression.ConditionBuilder {
	// It's possible to use ISO8601 string format with Geater than cf https://www.abhayachauhan.com/2017/12/how-to-store-dates-or-timestamps-in-dynamodb/
	// RFC3339 is some standard based on and stricter than ISO8601
	return expression.GreaterThan(expression.Name("created_at"), expression.Value(now.BeginningOfQuarter().Format(time.RFC3339)))
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
