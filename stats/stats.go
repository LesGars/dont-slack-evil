package stats

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

var thresholdQuality float64 = 0.33 // Means negative over 33% --> bad quality
var thresholdAlert float64 = 50.0

type DSEHomeStats struct {
	MessagesAnalyzedAllTime                  int     `json:"messagesAnalyzedAllTime"`
	MessagesOfGoodQualityAllTime             int     `json:"messagesOfBadQualityAllTime"`
	PercentageOfMessagesOfGoodQualityAllTime float64 `json:"percentageOfMessagesOfBadQualityAllTime"`

	MessagesAnalyzedSinceQuarter                  int     `json:"messagesAnalyzedLastQuarter"`
	MessagesOfGoodQualitySinceQuarter             int     `json:"messagesOfBadQualityLastQuarter"`
	PercentageOfMessagesOfGoodQualitySinceQuarter float64 `json:"percentageOfMessagesOfBadQualityLastQuarter"`
}

func HomeStatsForUser(userId string) DSEHomeStats {
	userIdFilt := userIdFilt(userId)
	goodQlFilt := goodQualityFilt()
	sinceBeginningOfQuarterFilt := sinceBeginningOfQuarterFilt()
	stats := DSEHomeStats{
		MessagesAnalyzedAllTime:      messagesAnalyzed(userIdFilt),
		MessagesOfGoodQualityAllTime: messagesAnalyzed(expression.And(goodQlFilt, userIdFilt)),

		MessagesAnalyzedSinceQuarter:      messagesAnalyzed(expression.And(userIdFilt, sinceBeginningOfQuarterFilt)),
		MessagesOfGoodQualitySinceQuarter: messagesAnalyzed(expression.And(goodQlFilt, userIdFilt, sinceBeginningOfQuarterFilt)),
	}
	if stats.MessagesAnalyzedAllTime != 0 {
		stats.PercentageOfMessagesOfGoodQualityAllTime = float64(stats.MessagesOfGoodQualityAllTime) / float64(stats.MessagesAnalyzedAllTime)
	}
	if stats.MessagesAnalyzedSinceQuarter != 0 {
		stats.PercentageOfMessagesOfGoodQualitySinceQuarter = float64(stats.MessagesOfGoodQualitySinceQuarter) / float64(stats.MessagesAnalyzedSinceQuarter)
	}
	return stats
}

func userIdFilt(userId string) expression.ConditionBuilder {
	return expression.Equal(expression.Name("user_id"), expression.Value(userId))
}

func goodQualityFilt() expression.ConditionBuilder {
	return expression.LessThan(expression.Name("sentiment.negative"), expression.Value(thresholdQuality))
}

func fromLastWeekFilt() expression.ConditionBuilder {
	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)
	return expression.GreaterThan(expression.Name("created_at"), expression.Value(lastWeek.Format(time.RFC3339)))
}

func sinceBeginningOfQuarterFilt() expression.ConditionBuilder {
	// It's possible to use ISO8601 string format with Geater than cf https://www.abhayachauhan.com/2017/12/how-to-store-dates-or-timestamps-in-dynamodb/
	// RFC3339 is some standard based on and stricter than ISO8601
	return expression.GreaterThan(expression.Name("created_at"), expression.Value(now.BeginningOfQuarter().Format(time.RFC3339)))
}

// GetWeeklyStats gets the weekly positivity score of a user
func GetWeeklyStats(userID string) (int, int) {
	userIDFilt := userIdFilt(userID)
	goodQualityFilt := goodQualityFilt()
	lastWeekFilt := fromLastWeekFilt()
	badMessages := messagesAnalyzed(expression.And(goodQualityFilt, userIDFilt, lastWeekFilt))
	totalMessages := messagesAnalyzed(expression.And(userIDFilt, lastWeekFilt))
	goodMessages := totalMessages - badMessages
	return goodMessages, totalMessages
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

// HasTooManyBadQualityMessagesLastQuarter returns true if the user sent too many messages of bad quality
// over the last quarter...
// TODO replace the stat by PercentageOfMessagesOfBadQualityLastQuarter
func HasTooManyBadQualityMessagesLastQuarter(userId string) bool {
	userStats := HomeStatsForUser(userId)
	return (userStats.PercentageOfMessagesOfGoodQualitySinceQuarter)*100 <= thresholdAlert
}
