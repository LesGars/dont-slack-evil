package apphome

import (
	"testing"

	"github.com/go-test/deep"
)

func TestHomeStatsForUser(t *testing.T) {
	// Stats correspond to arbitrary messages that were seeded in the test DB
	// Using https://snippets.cacher.io/snippet/3ae7dcb2e44370bf4dfc
	expectedObject := DSEHomeStats{50, 29, 0.58, 26, 14, 0.5384615384615384}
	actual := HomeStatsForUser("Mr. Maggie Feest")

	if diff := deep.Equal(expectedObject, actual); diff != nil {
		t.Error(diff)
	}
}
