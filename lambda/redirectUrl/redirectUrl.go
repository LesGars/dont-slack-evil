package main

import (
	"log"
	"os"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	Entreprise string `json:"entreprise"`
}
// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
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
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%s", body)
	var oauthAccessResponse OauthAccessResponse;
	unMarshallErr := json.Unmarshal(body, &oauthAccessResponse)
	if unMarshallErr != nil {
		log.Println(unMarshallErr)
	}
	log.Println(oauthAccessResponse)
	botAccessToken := oauthAccessResponse.AccessToken
	log.Println(botAccessToken)
	redirectURL := "https://app.slack.com/client/" + oauthAccessResponse.Team.ID
	response := Response{
		StatusCode: 302,
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
