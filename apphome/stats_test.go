package apphome

import (
	"testing"

	"github.com/go-test/deep"
)

func TestHomeStatsForUser(t *testing.T) {
	// Stats correspond to arbitrary messages that were seeded in the test DB
	// Using https://snippets.cacher.io/snippet/3ae7dcb2e44370bf4dfc
	expectedObject := DSEHomeStats{50, 27, 0.54, 0, 0, 0}
	actual := HomeStatsForUser("Alissa Kutch")

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
