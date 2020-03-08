package db

type Message struct {
	UserId         string    `json:"user_id"`
	SlackMessageId string    `json:"slack_message_id"`
	SlackThreadId  string    `json:"slack_thread_id"`
	SlackTeamId    string    `json:"slack_team_id"`
	Text           string    `json:"text"`
	Analyzed       bool      `json:"analyzed"`
	CreatedAt      string    `json:"created_at"`
	Quality        float64   `json:"quality"`
	Sentiment      Sentiment `json:"sentiment"`
}

type Sentiment struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}
