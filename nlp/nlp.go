package nlp

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SentimentAnalysis is the response type of the GetSentiment func
type SentimentAnalysis struct {
  Sentiment string `json:"sentiment"`
  Message string `json:"message"`
}

// GetSentiment computes a percentage of happy/neutral/sad for a given string
func GetSentiment(message string, apiURL string, apiKey string) (SentimentAnalysis, error) {
	// Get sentiment of message
	form := url.Values{}
	form.Add("text", message)
	form.Add("api_key", apiKey)
	resp, err := http.Post(
		apiURL + "/v4/sentiment",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if (err != nil) {
		return SentimentAnalysis{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	responseBody := SentimentAnalysis{
		Sentiment: string(body),
		Message: message,
	}

	return responseBody, nil
}
