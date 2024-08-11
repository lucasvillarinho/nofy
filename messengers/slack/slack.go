package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lucasvillarinho/nofy"
)

const Timeout = 5000

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Slack is a client to send messages to Slack.
type Slack struct {
	URL     string
	Token   string
	Timeout time.Duration
	Client  HTTPClient
	Message Message
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
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Option func(*Slack)

// NewSlackClient creates a new Slack client.
func NewSlackMensseger(options ...Option) (nofy.Messenger, error) {
	slack := &Slack{
		URL:     "https://slack.com/api/chat.postMessage",
		Timeout: Timeout * time.Millisecond,
	}

	for _, opt := range options {
		opt(slack)
	}

	if len(strings.TrimSpace(slack.Token)) == 0 {
		return nil, fmt.Errorf("missing token")
	}
	if slack.Timeout == 0 {
		return nil, fmt.Errorf("missing timeout")
	}
	if len(strings.TrimSpace(slack.Message.Channel)) == 0 {
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

// Send sends a message to a Slack channel.
// It returns an error if the message could not be sent,
// or if the response from Slack is not OK.
func (s *Slack) sendRequest(ctx context.Context, body []byte) (*Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.URL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := s.Client
	if client == nil {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"error sending message. Status Code: %d",
			resp.StatusCode,
		)
	}

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var slackResponse Response
	err = json.Unmarshal(bodyResponse, &slackResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &slackResponse, nil
}

// Send sends a message with blocks to a Slack channel.
// Block messages are used to create rich messages with elements.
// Doc: https://api.slack.com/reference/messaging/blocks
// Playground: https://app.slack.com/block-kit-builder
func (s *Slack) Send(ctx context.Context) error {
	jsonMessage, err := json.Marshal(s.Message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	slackResponse, err := s.sendRequest(ctx, jsonMessage)
	if err != nil {
		return err
	}

	if !slackResponse.OK {
		return fmt.Errorf(
			"error sending message: %s",
			slackResponse.Error,
		)
	}
	return nil
}
