package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pool "github.com/alitto/pond"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/helpers"
)

const Timeout = 5000

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Slack is a client to send messages to Slack.
type Slack struct {
	ID         string
	URL        string
	Token      string
	Timeout    time.Duration
	Recipients []string
	Client     HTTPClient
	Message    []map[string]any
}

// Response is the response from Slack.
// OK is true if the message was sent successfully.
// Error contains the error message if the message could not be sent.
type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// BlockMessage is a message with blocks to send to Slack.
type BlockMessage struct {
	Channel string           `json:"channel"`
	Blocks  []map[string]any `json:"blocks"`
}
type Option func(*Slack)

// NewSlackClient creates a new Slack client.
func NewSlackClient(options ...Option) (nofy.Sender, error) {
	slack := &Slack{
		ID:         "slack-sender-" + helpers.GenerateUUID(),
		URL:        "https://slack.com/api/chat.postMessage",
		Timeout:    Timeout * time.Millisecond,
		Recipients: make([]string, 0),
	}

	for _, opt := range options {
		opt(slack)
	}

	if slack.Token == "" {
		return nil, fmt.Errorf("missing Slack Token")
	}
	if slack.Timeout == 0 {
		return nil, fmt.Errorf("missing Timeout")
	}

	return slack, nil
}

// WithToken sets the Token for the Slack client.
func WithToken(Token string) Option {
	return func(s *Slack) {
		s.Token = Token
	}
}

// WithTimeout sets the Timeout for the Slack client.
func WithTimeout(Timeout time.Duration) Option {
	return func(s *Slack) {
		s.Timeout = Timeout
	}
}

func (s *Slack) GetId() string {
	return s.ID
}

// RemoveRecipient removes a recipient from the Slack client.
func (s *Slack) RemoveRecipient(recipient any) error {
	recipient, ok := recipient.(string)
	if !ok {
		return fmt.Errorf("invalid recipient")
	}

	for i, r := range s.Recipients {
		if r == recipient {
			s.Recipients = append(s.Recipients[:i], s.Recipients[i+1:]...)
			return nil
		}
	}
	return nil
}

// AddRecipient adds a recipient to the Slack client.
func (s *Slack) AddRecipient(recipient any) error {
	r, ok := recipient.(string)
	if !ok {
		return fmt.Errorf("invalid recipient")
	}
	s.Recipients = append(s.Recipients, r)
	return nil
}

// Send sends a message to a Slack channel.
// It returns an error if the message could not be sent,
// or if the response from Slack is not OK.
func (s *Slack) send(ctx context.Context, body []byte) (*Response, error) {
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
		return nil, fmt.Errorf("error sending message: %s", resp.Status)
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
// Block messages are used to create rich messages with buttons, images, and other elements.
// Doc https://api.slack.com/reference/messaging/blocks
// Playground https://app.slack.com/block-kit-builder
func (s *Slack) Send(ctx context.Context) error {
	pool := pool.New(len(s.Recipients), len(s.Recipients))
	group, ctx := pool.GroupContext(ctx)

	for _, channel := range s.Recipients {
		channel := channel
		group.Submit(func() error {
			message := BlockMessage{
				Channel: channel,
				Blocks:  s.Message,
			}
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("error marshalling message: %w", err)
			}

			slackResponse, err := s.send(ctx, jsonMessage)
			if err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}

			if !slackResponse.OK {
				return fmt.Errorf(
					"error sending message: %s",
					slackResponse.Error,
				)
			}
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}
