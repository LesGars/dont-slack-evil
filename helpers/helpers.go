package helpers

import (
	"regexp"

	"rsc.io/quote"
)

// Hello returns "hello world" translated in the system's locale
func Hello() string {
	return quote.Hello()
}

// ReverseRunes returns its argument string reversed rune-wise left to right.
func ReverseRunes(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func QuoteForSlack(message string) string {
	var re = regexp.MustCompile(`(.+)`)
	return re.ReplaceAllString(message, `> $1`)
}
