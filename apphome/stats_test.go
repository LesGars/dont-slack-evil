package apphome

import (
	dsedb "dont-slack-evil/db"
	"log"
	"strconv"
	"testing"

	"github.com/fatih/structs"
	"github.com/go-test/deep"
)

func TestHomeStatsForUser(t *testing.T) {
	for count := 0; count <= 2; count++ {
		msg := dsedb.Message{
			UserId:         "42",
			SlackMessageId: strconv.Itoa(count),
		}
		log.Printf("Next step is to send %s to DYnamoDB", structs.Values(msg))
		// dsedb.Store(
		// 	os.Getenv("DYNAMODB_TABLE"),
		// 	structs.Map(&msg),
		// )
	}

	expectedObject := DSEHomeStats{3, 0, 0}
	actual := HomeStatsForUser("42")

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		// t.Error(diff)
	}
}
