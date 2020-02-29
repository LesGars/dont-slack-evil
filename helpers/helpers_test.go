package helpers

import (
	"testing"

	"github.com/go-test/deep"
)

func TestReverseRunes(t *testing.T) {
	want := "olleH"
	if got := ReverseRunes("Hello"); got != want {
		t.Errorf("ReverseRunes() = %q, want %q", got, want)
	}
}

func TestQuoteForSlack(t *testing.T) {
	expected := "> Hello\n> How are you?\n\n> Life is lemons!"
	actual := QuoteForSlack("Hello\nHow are you?\n\nLife is lemons!")

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}
