package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// If parseEvent fails, the handler should return an error
func TestHandleEvent_ParseEventFailure(t *testing.T) {
	old := ParseEvent
	defer func() { ParseEvent = old }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{}, errors.New("Error-Mock")
	}
	_, e := HandleEvent([]byte("abcd"))
	want := "Could not parse Slack event :'("
	got := e.Error()
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type URLVerification and the JSONification fails, the handler should return an error
func TestHandleEvent_URLVerificationFailure(t *testing.T) {
	old := ParseEvent
	defer func() { ParseEvent = old }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{Type: slackevents.URLVerification}, nil
	}
	_, e := HandleEvent([]byte("{{}"))
	want := "Unable to register the URL"
	got := e.Error()
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}
}

// If parseEvent returns an event of type URLVerification and the JSONification does not fail, the handler should return 
// the challenge value and nil
func TestHandleEvent_URLVerificationSuccess(t *testing.T) {
	old := ParseEvent
	defer func() { ParseEvent = old }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{Type: slackevents.URLVerification}, nil
	}
	got, e := HandleEvent([]byte("{\"Challenge\": \"Challenge\"}"))
	want := "Challenge"
	
	if e != nil {
		t.Errorf("The handler unexpectedly returned an error")
	}
	if got != want {
		t.Errorf("The handler doesn't return the right challenge, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type AppMentionEvent and the POST message fails, the handler should return an error
func TestHandleEvent_AppMentionEventFailure(t *testing.T) {
	oldParseEvent := ParseEvent
	defer func() { ParseEvent = oldParseEvent }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type: slackevents.AppMention, 
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppMentionEvent{}},
			}, nil
	}

	oldPostMessage := PostMessage
	defer func() { PostMessage = oldPostMessage }()
	PostMessage = func(channelID string, options ...slack.MsgOption) (string, string, error) {
		return "", "", errors.New("Error-mock")
	}
	_, e := HandleEvent([]byte("{\"Challenge\": \"Challenge\"}"))
	want := "Error while posting message Error-mock"
	got := e.Error()
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type AppMentionEvent and the POST message succeeds, the handler should return nil and nil
func TestHandleEvent_AppMentionEventSuccess(t *testing.T) {
	oldParseEvent := ParseEvent
	defer func() { ParseEvent = oldParseEvent }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type: slackevents.AppMention, 
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppMentionEvent{}},
			}, nil
	}

	oldPostMessage := PostMessage
	defer func() { PostMessage = oldPostMessage }()
	PostMessage = func(channelID string, options ...slack.MsgOption) (string, string, error) {
		return "", "", nil
	}
	resp, e := HandleEvent([]byte("{\"Challenge\": \"Challenge\"}"))

	if resp != "" {
		t.Errorf("The handler should have returned no response. Instead it returned %v", resp)
	}
	if e != nil {
		t.Errorf("The handler should not failed. It returned the following error %v", e)
	}

}

// If parseEvent returns an event of type AppHomeOpened and the POST message fails, the handler should return an error
func TestHandleEvent_AppHomeOpenedFailure(t *testing.T) {
	oldParseEvent := ParseEvent
	defer func() { ParseEvent = oldParseEvent }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type: slackevents.CallbackEvent, 
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppHomeOpenedEvent{}},
			}, nil
	}

	oldPublishView := PublishView
	defer func() { PublishView = oldPublishView }()
	publishViewError := errors.New("Error-Mock")
	PublishView = func(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error) {
		return nil, publishViewError
	}
	_, got := HandleEvent([]byte("{\"Challenge\": \"Challenge\"}"))
	want := publishViewError
	if got != want {
		t.Errorf("The handler doesn't return the right error, got %v want %v", got, want)
	}

}

// If parseEvent returns an event of type AppHomeOpened and the POST message succeeds, the handler should return nil-nil
func TestHandleEvent_AppHomeOpenedSuccess(t *testing.T) {
	oldParseEvent := ParseEvent
	defer func() { ParseEvent = oldParseEvent }()
	ParseEvent = func(rawEvent json.RawMessage, opts ...slackevents.Option) (slackevents.EventsAPIEvent, error) {
		return slackevents.EventsAPIEvent{
			Type: slackevents.CallbackEvent, 
			InnerEvent: slackevents.EventsAPIInnerEvent{Data: &slackevents.AppHomeOpenedEvent{}},
			}, nil
	}

	oldPublishView := PublishView
	defer func() { PublishView = oldPublishView }()

	PublishView = func(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error) {
		return nil, nil
	}
	resp, e := HandleEvent([]byte("{\"Challenge\": \"Challenge\"}"))

	if resp != "" {
		t.Errorf("The handler should have returned no response. Instead it returned %v", resp)
	}
	if e != nil {
		t.Errorf("The handler should not failed. It returned the following error %v", e)
	}

}