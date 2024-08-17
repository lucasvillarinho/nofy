package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/helpers/request"
)

const Timeout = 5000

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Slack is a client to send messages to Slack.
type Slack struct {
	URL       string
	Token     string
	Message   Message
	Timeout   time.Duration
	requester request.Requester
}

// Message is the message to send to Slack.
type Message struct {
	Channel string           `json:"channel"`
	Content []map[string]any `json:"blocks"`
}

// Response is the response from Slack.
// OK is true if the message was sent successfully.
// Error contains the error message if the message could not be sent.
type Response struct {
	Error string `json:"error,omitempty"`
	OK    bool   `json:"ok"`
}

type Option func(*Slack)

// NewSlackMessenger creates a new Slack client.
func NewSlackMessenger(options ...Option) (nofy.Messenger, error) {
	slack := &Slack{
		URL:     "https://slack.com/api/chat.postMessage",
		Timeout: Timeout * time.Millisecond,
	}

	for _, opt := range options {
		opt(slack)
	}

	if strings.TrimSpace(slack.Token) == "" {
		return nil, fmt.Errorf("missing token")
	}
	if slack.Timeout == 0 {
		return nil, fmt.Errorf("missing timeout")
	}
	if strings.TrimSpace(slack.Message.Channel) == "" {
		return nil, fmt.Errorf("missing channel")
	}
	if slack.Message.Content == nil {
		return nil, fmt.Errorf("missing content")
	}

	return slack, nil
}

// WithToken sets the Token for the Slack client.
func WithToken(token string) Option {
	return func(s *Slack) {
		s.Token = token
	}
}

// WithTimeout sets the Timeout for the Slack client.
func WithTimeout(timeout time.Duration) Option {
	return func(s *Slack) {
		s.Timeout = timeout
	}
}

// WithMessage sets the Message for the Slack client.
func WithMessage(message Message) Option {
	return func(s *Slack) {
		s.Message = message
	}
}

// Send sends a message with blocks to a Slack channel.
// Block messages are used to create rich messages with elements.
// Doc: https://api.slack.com/reference/messaging/blocks
// Playground: https://app.slack.com/block-kit-builder
func (s *Slack) Send(ctx context.Context) error {
	msg, err := json.Marshal(s.Message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	httpClient := &http.Client{
		Timeout: s.Timeout,
	}

	resp, body, err := s.requester.DoWithCtx(ctx,
		request.WithMethod(http.MethodPost),
		request.WithURL(s.URL),
		request.WithPayload(msg),
		request.WithClient(httpClient))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"error sending message status-code: %d",
			resp.StatusCode,
		)
	}

	var slackResponse Response
	err = json.Unmarshal(body, &slackResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	if !slackResponse.OK {
		return fmt.Errorf(
			"error sending message: %s",
			slackResponse.Error,
		)
	}
	return nil
}
