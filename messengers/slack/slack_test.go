package slack

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/lcvilla/nofy/helpers/assert"
	"github.com/lcvilla/nofy/helpers/request"
)

func TestNewSlackMessenger(t *testing.T) {
	t.Run("should create Slack messenger successfully", func(t *testing.T) {
		messenger, err := NewSlackMessenger(
			WithToken("test-token"),
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
		assert.IsNil(t, err)
		assert.IsNotNil(t, messenger)
	})

	t.Run("should return error when token is missing", func(t *testing.T) {
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

	t.Run("should return error when timeout is invalid", func(t *testing.T) {
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

	t.Run("should return error when channel is missing", func(t *testing.T) {
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

	t.Run("should return error when message content is missing", func(t *testing.T) {
		_, err := NewSlackMessenger(
			WithToken("test-token"),
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					Channel: "test-channel",
				}),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing message"),
			"Expected missing message error",
		)
	})
}

func TestSlackOptions(t *testing.T) {
	t.Run("should set token correctly with WithToken option", func(t *testing.T) {
		slack := &Slack{}
		WithToken("test-token")(slack)

		assert.AreEqual(
			t,
			slack.Token,
			"test-token",
			"Expected token to be 'test-token'",
		)
	})

	t.Run("should set timeout correctly with WithTimeout option", func(t *testing.T) {
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
	t.Run("should send message successfully", func(t *testing.T) {
		mockRequester := &request.MockRequester{
			DoFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
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

	t.Run("should return error when marshalling message fails", func(t *testing.T) {
		msg := Message{
			Channel: "test-channel",
			Content: []map[string]any{
				{
					"type": "section",
					"text": map[string]any{
						"type": make(chan string),
						"text": "Hello, World!",
					},
				},
			},
		}
		messenger := &Slack{
			Message:   msg,
			URL:       "https://slack.com/api/chat.postMessage",
			Timeout:   5 * time.Second,
			requester: nil,
		}

		err := messenger.Send(context.TODO())

		assert.AreEqualErrs(
			t,
			err,
			errors.New("error marshaling message: json: unsupported type: chan string"),
			"Expected error marshalling message",
		)
	})

	t.Run("should return error when status code is not OK", func(t *testing.T) {
		mockRequester := &request.MockRequester{
			DoFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
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
