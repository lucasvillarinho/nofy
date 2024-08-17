package slack

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/lucasvillarinho/nofy/helpers/assert"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type FailingReader struct{}

func (r *FailingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading response")
}

func TestNewSlackMessenger(t *testing.T) {
	t.Run("missing Token", func(t *testing.T) {
		_, err := NewSlackMessenger(
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					Channel: "test-channel",
					Content: []map[string]any{
						{
							"type": "section",
							"text": map[string]string{
								"type": "mrkdwn",
								"text": "Hello, World!",
							},
						},
					},
				}),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing token"),
			"Expected missing Slack Token error",
		)
	})

	t.Run("invalid Timeout", func(t *testing.T) {
		_, err := NewSlackMessenger(
			WithToken("test-token"),
			WithTimeout(0),
			WithMessage(
				Message{
					Channel: "test-channel",
					Content: []map[string]any{
						{
							"type": "section",
							"text": map[string]string{
								"type": "mrkdwn",
								"text": "Hello, World!",
							},
						},
					},
				}),
		)

		assert.AreEqual(
			t,
			err,
			errors.New("missing timeout"),
			"Expected missing Timeout error",
		)
	})

	t.Run("missing channel", func(t *testing.T) {
		_, err := NewSlackMessenger(
			WithToken("test-token"),
			WithTimeout(5*time.Second),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing channel"),
			"Expected missing Message error",
		)
	})
}

func TestSlackOptions(t *testing.T) {
	t.Run("withToken", func(t *testing.T) {
		slack := &Slack{}
		WithToken("test-token")(slack)

		assert.AreEqual(
			t,
			slack.Token,
			"test-token",
			"Expected token to be 'test-token'",
		)
	})

	t.Run("withTimeout", func(t *testing.T) {
		slack := &Slack{}
		WithTimeout(10 * time.Second)(slack)

		assert.AreEqual(
			t,
			slack.Timeout,
			10*time.Second,
			"Expected timeout to be 10s",
		)
	})
}
