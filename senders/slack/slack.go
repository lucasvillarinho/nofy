package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	envl "github.com/caarlos0/env/v11"
)

const timeout = 5000

// Slack is a client to send messages to Slack.
type Slack struct {
	url        string
	token      string
	timeout    time.Duration
	recipients []Recipient
}

// config is the configuration for the Slack client.
type config struct {
	token string `env:"NOFY_SLACK_TOKEN"`
}

// Message is the message to send to Slack.
type Message struct {
	Channel  string `json:"channel"`
	Text     string `json:"text"`
	Markdown bool   `json:"mrkdwn"`
}

// Response is the response from Slack.
// OK is true if the message was sent successfully.
// Error contains the error message if the message could not be sent.
type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Recipient struct {
	Channel string `json:"channel"`
}

type BlockMessage struct {
	Channel string           `json:"channel"`
	Blocks  []map[string]any `json:"blocks"`
}
type Option func(*Slack)

// NewSlackClient creates a new Slack client.
func NewSlackClient(options ...Option) (*Slack, error) {
	cfg := config{}
	if err := envl.Parse(&cfg); err != nil {
		return nil, err
	}
	slack := &Slack{
		url:        "https://slack.com/api/chat.postMessage",
		token:      cfg.token,
		timeout:    timeout * time.Millisecond,
		recipients: make([]Recipient, 0),
	}

	for _, opt := range options {
		opt(slack)
	}

	return slack, nil
}

// WithTimeout sets the timeout for the Slack client.
func WithTimeout(timeout time.Duration) Option {
	return func(s *Slack) {
		s.timeout = timeout
	}
}

// WithToken sets the token for the Slack client.
func WithToken(token string) Option {
	return func(s *Slack) {
		s.token = token
	}
}

// AddRecipient adds a recipient to the list of recipients.
func (s *Slack) AddRecipient(channel string) {
	s.recipients = append(s.recipients, Recipient{Channel: channel})
}

// RemoveRecipient removes a recipient from the list of recipients.
func (s *Slack) RemoveRecipient(channel string) {
	for i, recipient := range s.recipients {
		if recipient.Channel == channel {
			s.recipients = append(s.recipients[:i], s.recipients[i+1:]...)
			return
		}
	}
}

// Send sends a message to a Slack channel.
// It returns an error if the message could not be sent,
// or if the response from Slack is not OK.
func (s *Slack) send(ctx context.Context, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	defer resp.Body.Close()

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var slackResponse Response
	err = json.Unmarshal(bodyResponse, &slackResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	if !slackResponse.OK {
		return fmt.Errorf("error sending message: %s", slackResponse.Error)
	}

	return nil
}

// SendBlocks sends a message with blocks to a Slack channel.
// Block messages are used to create rich messages with buttons, images, and other elements.
// Doc https://api.slack.com/reference/messaging/blocks
// Playground https://app.slack.com/block-kit-builder
func (s *Slack) SendBlocks(ctx context.Context, blocks []map[string]any) error {
	for _, re := range s.recipients {
		message := BlockMessage{
			Channel: re.Channel,
			Blocks:  blocks,
		}
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("error marshalling message: %w", err)
		}

		err = s.send(ctx, jsonMessage)
		if err != nil {
			return fmt.Errorf("error sending message with blocks: %w", err)
		}
	}

	return nil
}

// Send sends a message to a Slack channel.
// It returns an error if the message could not be sent,
// or if the response from Slack is not OK.
// The message is sent to all the channels in the list.
func (s *Slack) Send(ctx context.Context, msg string) error {
	for _, re := range s.recipients {
		msg := Message{
			Channel:  re.Channel,
			Text:     msg,
			Markdown: true,
		}

		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("error marshalling message: %w", err)
		}

		err = s.send(ctx, jsonMessage)
		if err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}
	}
	return nil
}
