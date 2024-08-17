package slack

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/lucasvillarinho/nofy/helpers/assert"
	"github.com/lucasvillarinho/nofy/helpers/request"
)

// errorReader é uma estrutura que implementa io.Reader, mas sempre retorna um erro
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

// MockRequest é uma estrutura que você pode usar para mockar as requisições HTTP
type MockRequest struct {
	DoFunc func(ctx context.Context, options ...request.Option) (*http.Response, error)
}

func (m *MockRequest) DoWithCtx(
	ctx context.Context,
	options ...request.Option,
) (*http.Response, error) {
	return m.DoFunc(ctx, options...)
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

func TestSlackSend(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		mockRequester := &request.MockRequester{
			DoWithCtxFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, []byte(`{"ok": true}`), nil
			},
		}

		msg := Message{
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
		}
		messenger := &Slack{
			Message:   msg,
			URL:       "https://slack.com/api/chat.postMessage",
			Timeout:   5 * time.Second,
			requester: mockRequester,
		}

		err := messenger.Send(context.TODO())

		assert.IsNil(t, err)
	})

	t.Run("status nok", func(t *testing.T) {
		mockRequester := &request.MockRequester{
			DoWithCtxFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
				}, []byte(`{"ok": false}`), nil
			},
		}

		msg := Message{
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
		}
		messenger := &Slack{
			Message:   msg,
			URL:       "",
			Timeout:   5 * time.Second,
			requester: mockRequester,
		}

		err := messenger.Send(context.TODO())

		assert.AreEqualErrs(
			t,
			err,
			errors.New("error sending message status-code: 400"),
			"Expected error sending message",
		)
	})
}
