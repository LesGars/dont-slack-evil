package test

import (
	dsedb "dont-slack-evil/db"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func MockWorkingApiForTeam(slackevents.EventsAPIEvent) (*dsedb.ApiForTeam, error) {
	workingApi := WorkingApiForTeam()
	return &workingApi, nil
}

func WorkingApiForTeam() dsedb.ApiForTeam {
	return dsedb.ApiForTeam{
		Team:                  dsedb.Team{SlackTeamId: "42", SlackBotUserToken: "xoxb-42"},
		SlackBotUserApiClient: WorkingDummySlackClient{},
	}
}

type WorkingDummySlackClient struct{}

func (ds WorkingDummySlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	return "", "", nil
}
func (ds WorkingDummySlackClient) PublishView(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error) {
	return &slack.ViewResponse{}, nil
}
func (ds WorkingDummySlackClient) GetUserInfo(user string) (*slack.User, error) {
	return &slack.User{Name: "Le gars"}, nil
}
func (ds WorkingDummySlackClient) GetUsers() ([]slack.User, error) {
	return []slack.User{
		slack.User{Name: "Le gars", ID: "42"},
		slack.User{Name: "La meuf", ID: "44"},
		slack.User{Name: "L'autre'", ID: "22"},
	}, nil
}
