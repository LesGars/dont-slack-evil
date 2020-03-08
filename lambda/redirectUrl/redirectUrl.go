package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	dsedb "dont-slack-evil/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fatih/structs"
)

// OauthAccessResponse is the type of the Oauth Access response
type OauthAccessResponse struct {
	Ok bool `json:"ok"`
	AppID string `json:"app_id"`
	Scope string `json:"scope"`
	TokenType string `json:"token_type"`
	AccessToken string `json:"access_token"`
	BotUserID string `json:"bot_user_id"`
	Team struct {
		ID string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	Enterprise string `json:"enterprise"`
}

// OauthTokenDBItem is the struct for storing the access token in DB
type OauthTokenDBItem struct {
	TeamID string `json:"team_id"`
	AccessToken string `json:"access_token"`
}
// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	structs.DefaultTagName = "json" // https://github.com/fatih/structs/issues/25
	var statusCode int;
	statusCode = 302;
	// Get Oauth Access token
	query := request.QueryStringParameters
	code := query["code"]
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	oauthURL := "https://slack.com/api/oauth.v2.access"
	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	resp, err := http.Post(
		oauthURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		log.Println(err)
		statusCode = 500
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var oauthAccessResponse OauthAccessResponse;
	unMarshallErr := json.Unmarshal(body, &oauthAccessResponse)
	if unMarshallErr != nil {
		log.Println(unMarshallErr)
		statusCode = 500
	}

	// Save in DB
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "teams"
	dbError := dsedb.CreateTableIfNotCreated(tableName, "slack_team_id")
	if dbError {
		log.Println(dbError)
		statusCode = 500
	}
	dbItem := dsedb.Team{
		SlackTeamId: oauthAccessResponse.Team.ID,
		SlackBotUserToken: oauthAccessResponse.AccessToken,
	}
	dbResult := dsedb.Store(tableName, structs.Map(&dbItem))
	if !dbResult {
		log.Println("Could not store message in DB")
		statusCode = 500
	} else {
		log.Println("Oauth Access token was stored successfully")
	}

	// Redirect to slack workspace URL
	redirectURL := "https://app.slack.com/client/" + oauthAccessResponse.Team.ID
	response := Response{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Location": redirectURL,
		},
		Body: "",
	}
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
