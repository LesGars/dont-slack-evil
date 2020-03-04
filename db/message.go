package db

type Message struct {
	Id             string  `json:"id"`
	UserId         string  `json:"user_id"`
	SlackMessageId string  `json:"slack_message_id"`
	Analyzed       bool    `json:"analyzed"`
	CreatedAt      string  `json:"created_at"`
	Quality        float64 `json:"quality"`
	Sentiment      Sentiment
}

type Sentiment struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}
