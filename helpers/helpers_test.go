package helpers

import "testing"

func TestReverseRunes(t *testing.T) {
	want := "olleH"
	if got := ReverseRunes("Hello"); got != want {
		t.Errorf("ReverseRunes() = %q, want %q", got, want)
	}
}
