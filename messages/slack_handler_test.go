package messages

import (
	"dont-slack-evil/db"
	"encoding/json"
	"errors"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// If parseEvent fails, the handler should return an error
func TestSlackHandler_parseEventFailure(t *testing.T) {
	old := parseEvent
	defer func() { parseEvent = old }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{}, errors.New("Error-Mock")
	}
	_, e := SlackHandler([]byte("abcd"))
	want := "Could not parse Slack event :'("
	got := e.Error()
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type URLVerification and the JSONification fails, the handler should return an error
func TestSlackHandler_URLVerificationFailure(t *testing.T) {
	old := parseEvent
	defer func() { parseEvent = old }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{Type: slackevents.URLVerification}, nil
	}
	_, e := SlackHandler([]byte("{{}"))
	want := "Unable to register the URL"
	got := e.Error()
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}
}

// If parseEvent returns an event of type URLVerification and the JSONification does not fail, the handler should return
// the challenge value and nil
func TestSlackHandler_URLVerificationSuccess(t *testing.T) {
	old := parseEvent
	defer func() { parseEvent = old }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{Type: slackevents.URLVerification}, nil
	}
	got, e := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))
	want := "Challenge"

	if e != nil {
		t.Errorf("The handler unexpectedly returned an error")
	}
	if got != want {
		t.Errorf("The handler doesn't return the right challenge, got %v want %v", got, want)
	}
}

func mockApiForTeam() ApiForTeam {
	mockedSlackClient := slack.New("xoxb-42")

	return ApiForTeam{
		Team:                  db.Team{SlackTeamId: "42", SlackBotUserToken: "xoxb-42"},
		SlackBotUserApiClient: mockedSlackClient,
	}
}

/* Broken tests : see https://lesgarshack.slack.com/archives/CUHTQKV9N/p1583598765006100?thread_ts=1583597618.005900&cid=CUHTQKV9N

// If parseEvent returns an event of type AppMentionEvent and the POST message fails, the handler should return an error
func TestHandleSlackEvent_AppMentionEventFailure(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.AppMention,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppMentionEvent{}},
		}, nil
	}

	oldpostMessage := postMessage
	defer func() { postMessage = oldpostMessage }()
	postMessage = func(channelID string, options ...slack.MsgOption) (string, string, error) {
		return "", "", errors.New("Error-mock")
	}
	_, e := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))
	want := "Error while posting message Error-mock"
	got := e.Error()
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type AppMentionEvent and the POST message succeeds, the handler should return nil and nil
func TestHandleSlackEvent_AppMentionEventSuccess(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.AppMention,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppMentionEvent{}},
		}, nil
	}

	oldpostMessage := postMessage
	defer func() { postMessage = oldpostMessage }()
	postMessage = func(channelID string, options ...slack.MsgOption) (string, string, error) {
		return "", "", nil
	}
	resp, e := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))

	if resp != "" {
		t.Errorf("The handler should have returned no response. Instead it returned %v", resp)
	}
	if e != nil {
		t.Errorf("The handler should not failed. It returned the following error %v", e)
	}

}

If parseEvent returns an event of type AppHomeOpened and the POST message fails, the handler should return an error
func TestHandleSlackEvent_AppHomeOpenedFailure(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.CallbackEvent,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppHomeOpenedEvent{}},
		}, nil
	}

	olduserHome := userHome
	defer func() { userHome = olduserHome }()
	userHome = func(userId string) slack.Message {
		return slack.Message{}
	}

	oldpublishView := publishView
	defer func() { publishView = oldpublishView }()
	publishViewError := errors.New("Error-Mock")
	publishView = func(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error) {
		return nil, publishViewError
	}
	_, got := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))
	want := publishViewError
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}
}

If parseEvent returns an event of type AppHomeOpened and the POST message succeeds, the handler should return nil-nil
func TestHandleSlackEvent_AppHomeOpenedSuccess(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.CallbackEvent,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppHomeOpenedEvent{}},
		}, nil
	}

	oldpublishView := publishView
	defer func() { publishView = oldpublishView }()
	publishView = func(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error) {
		return nil, nil
	}

	olduserHome := userHome
	defer func() { userHome = olduserHome }()
	userHome = func(userId string) slack.Message {
		return slack.Message{}
	}

	resp, e := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))

	if resp != "" {
		t.Errorf("The handler should have returned no response. Instead it returned %v", resp)
	}
	if e != nil {
		t.Errorf("The handler should not failed. It returned the following error %v", e)
	}

}

// If parseEvent returns an event of type MessageEvent and storeMessage fails, the handler should return an error
func TestHandleSlackEvent_MessageEventStoreMessageFailure(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.CallbackEvent,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.MessageEvent{}},
		}, nil
	}

	oldstoreMessage := storeMessage
	defer func() { storeMessage = oldstoreMessage }()
	storeMessageError := errors.New("Error-Mock")
	storeMessage = func(message *slackevents.MessageEvent) error {
		return storeMessageError
	}
	_, got := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))
	want := storeMessageError
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type MessageEvent and getSentiment fails, the handler should return an error
func TestHandleSlackEvent_MessageEventGetSentimentFailure(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.CallbackEvent,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.MessageEvent{}},
		}, nil
	}

	oldstoreMessage := storeMessage
	defer func() { storeMessage = oldstoreMessage }()
	storeMessage = func(message *slackevents.MessageEvent) error {
		return nil
	}
	oldgetSentiment := getSentiment
	defer func() { getSentiment = oldgetSentiment }()
	getSentimentError := errors.New("Error-Mock")
	getSentiment = func(message *slackevents.MessageEvent) error {
		return getSentimentError
	}
	_, got := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))
	want := getSentimentError
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type MessageEvent and storeMessage/getSentiment succeeds, the handler should return nil-nil
func TestHandleSlackEvent_MessageEventSuccess(t *testing.T) {
	oldparseEvent := parseEvent
	defer func() { parseEvent = oldparseEvent }()
	parseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type:       slackevents.CallbackEvent,
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.MessageEvent{}},
		}, nil
	}

	oldstoreMessage := storeMessage
	defer func() { storeMessage = oldstoreMessage }()
	storeMessage = func(message *slackevents.MessageEvent) error {
		return nil
	}
	oldgetSentiment := getSentiment
	defer func() { getSentiment = oldgetSentiment }()
	getSentiment = func(message *slackevents.MessageEvent) error {
		return nil
	}
	resp, e := SlackHandler([]byte("{\"Challenge\": \"Challenge\"}"))

	if resp != "" {
		t.Errorf("The handler should have returned no response. Instead it returned %v", resp)
	}
	if e != nil {
		t.Errorf("The handler should not failed. It returned the following error %v", e)
	}
}
*/
