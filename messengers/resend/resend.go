package resend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lucasvillarinho/nofy/helpers/request"
)

const Timeout = 5000

var MarshalFunc = func(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Resend is a client to send messages to Resend.
type Resend struct {
	requester request.Requester
	URL       string
	Token     string
	Timeout   time.Duration
	Message   Message
}

// Message is the message to send to Resend.
// From is the email address of the sender (required).
// To is the email addresses of the recipients (required).
// CC is the email addresses of the CC recipients.
// Subject is the subject of the email (required).
// HTML is the HTML content of the email.
type Message struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	CC      string   `json:"cc"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

type Option func(*Resend)

// NewResendMessenger creates a new Resend client.
func NewResendMessenger(options ...Option) (*Resend, error) {
	resend := &Resend{
		URL:     "https://api.resend.com/emails",
		Timeout: Timeout * time.Millisecond,
	}

	for _, opt := range options {
		opt(resend)
	}

	err := validate(resend)
	if err != nil {
		return nil, err
	}

	return resend, nil
}

// validate validates the Resend client.
func validate(resend *Resend) error {
	if strings.TrimSpace(resend.Token) == "" {
		return fmt.Errorf("missing token")
	}
	if strings.TrimSpace(resend.Message.From) == "" {
		return fmt.Errorf("missing from")
	}
	if len(resend.Message.To) == 0 {
		return fmt.Errorf("missing to")
	}
	if strings.TrimSpace(resend.Message.Subject) == "" {
		return fmt.Errorf("missing subject")
	}

	return nil
}

// WithURL sets the URL for the Resend client.
func WithToken(token string) Option {
	return func(r *Resend) {
		r.Token = token
	}
}

// WithTimeout sets the Timeout for the Resend client.
func WithTimeout(timeout time.Duration) Option {
	return func(r *Resend) {
		r.Timeout = timeout
	}
}

// WithMessage sets the Message for the Resend client.
func WithMessage(message Message) Option {
	return func(r *Resend) {
		r.Message = message
	}
}

// Send sends a message using the Resend client.
func (r *Resend) Send(ctx context.Context) error {
	msg, err := MarshalFunc(r.Message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	HTTPClient := http.DefaultClient

	res, body, err := r.requester.Do(ctx,
		request.WithMethod(http.MethodPost),
		request.WithURL(r.URL),
		request.WithHeader("Authorization", "Bearer "+r.Token),
		request.WithHeader("Content-Type", "application/json"),
		request.WithHeader("Accept", "application/json"),
		request.WithClient(HTTPClient),
		request.WithPayload(msg),
	)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending message: status-code: %d body: %s", res.StatusCode, body)
	}

	return nil
}
