package slack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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

func TestNewSlackClient(t *testing.T) {
	t.Run("Missing Token", func(t *testing.T) {
		_, err := NewSlackMensseger(
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

	t.Run("Invalid Timeout", func(t *testing.T) {
		_, err := NewSlackMensseger(
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

	t.Run("Missing channel", func(t *testing.T) {
		_, err := NewSlackMensseger(
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

	t.Run("Successful client creation", func(t *testing.T) {
		slackClient, err := NewSlackMensseger(
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

		assert.AreEqual(
			t,
			slackClient.(*Slack).Token,
			"test-token",
			"Expected token to be 'test-token'",
		)
		assert.IsNil(t, err, "Expected timeout to be 5s")
	})
}

func TestSlackOptions(t *testing.T) {
	t.Run("WithToken option", func(t *testing.T) {
		slack := &Slack{}
		WithToken("test-token")(slack)

		assert.AreEqual(
			t,
			slack.Token,
			"test-token",
			"Expected token to be 'test-token'",
		)
	})

	t.Run("WithTimeout option", func(t *testing.T) {
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

func TestSlacksendRequest(t *testing.T) {
	t.Run("Successful message send", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				respBody := `{"ok": true}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString(respBody),
					),
					Header: make(http.Header),
				}, nil
			},
		}
		slack := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "test-token",
			Client: mockClient,
		}
		body := []byte(`{"text":"Hello, World!"}`)

		resp, err := slack.sendRequest(context.Background(), body)

		assert.IsNil(t, err, "Expected no error")
		assert.IsNotNil(t, resp, "Expected response to be not nil")
		assert.AreEqual(t, true, resp.OK, "Expected response to be OK")
	})

	t.Run("Error creating request", func(t *testing.T) {
		slack := &Slack{
			URL:    "://invalid-url",
			Token:  "test-token",
			Client: nil,
		}
		body := []byte(`{"text":"Hello, World!"}`)
		expectedErr := fmt.Errorf(
			"error creating request: %w",
			errors.New("parse \"://invalid-url\": missing protocol scheme"),
		)

		resp, err := slack.sendRequest(context.Background(), body)

		assert.IsNil(t, resp, "Expected nil response")
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error creating request",
		)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}
		slack := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "test-token",
			Client: mockClient,
		}
		body := []byte(`{"text":"Hello, World!"}`)
		expectedErr := fmt.Errorf(
			"error sending message: %w",
			errors.New("network error"),
		)

		resp, err := slack.sendRequest(context.Background(), body)

		assert.IsNil(t, resp, "Expected nil response")
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			expectedErr,
			"Expected error sending message",
		)
	})

	t.Run("Error due to non-OK status code", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     make(http.Header),
				}, nil
			},
		}
		slack := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "test-token",
			Client: mockClient,
		}
		body := []byte(`{"text":"Hello, World!"}`)
		errExpected := fmt.Errorf(
			"error sending message. Status Code: %d",
			http.StatusUnauthorized,
		)

		resp, err := slack.sendRequest(context.Background(), body)

		assert.IsNil(t, resp, "Expected nil response")
		assert.AreEqualErrs(t, err, errExpected, "Expected non-OK status code")
	})

	t.Run("Error reading response", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						&FailingReader{},
					),
					Header: make(http.Header),
				}, nil
			},
		}

		slack := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "test-token",
			Client: mockClient,
		}
		body := []byte(`{"text":"Hello, World!"}`)
		expectedErr := fmt.Errorf(
			"error reading response: %w",
			errors.New("error reading response"),
		)

		resp, err := slack.sendRequest(context.Background(), body)

		assert.IsNil(t, resp, "Expected nil response")
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error reading response",
		)
	})

	t.Run("Error unmarshalling response", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString("invalid json"),
					),
					Header: make(http.Header),
				}, nil
			},
		}
		slack := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "test-token",
			Client: mockClient,
		}
		expectedErr := fmt.Errorf(
			"error unmarshalling response: %w",
			errors.New("invalid character 'i' looking for beginning of value"),
		)
		body := []byte(`{"text":"Hello, World!"}`)

		resp, err := slack.sendRequest(context.Background(), body)

		assert.IsNil(t, resp, "Expected nil response")
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error unmarshalling response",
		)
	})
}

func TestSlackSend(t *testing.T) {
	t.Run("Successful message send", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				respBody := `{"ok": true}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respBody)),
					Header:     make(http.Header),
				}, nil
			},
		}
		slack := &Slack{
			Client: mockClient,
			Message: Message{
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
			},
		}

		err := slack.Send(context.Background())

		assert.IsNil(t, err, "Expected no error")
	})

	t.Run("Error marshaling message", func(t *testing.T) {
		slack := &Slack{
			Message: Message{
				Channel: "test-channel",
				Content: []map[string]any{
					{
						"type": func() {},
					},
				},
			},
		}
		errExpected := errors.New(
			"error marshaling message: json: unsupported type: func()",
		)

		err := slack.Send(context.Background())

		assert.AreEqualErrs(t, err, errExpected, "Expected marshaling error")
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}
		slack := &Slack{
			Client: mockClient,
			Message: Message{
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
			},
		}
		errExpected := errors.New("error sending message: network error")

		err := slack.Send(context.Background())

		assert.AreEqualErrs(
			t,
			err,
			errExpected,
			"Expected error sending message",
		)
	})

	t.Run("Error from Slack API", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				respBody := `{"ok": false, "error": "invalid_auth"}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respBody)),
					Header:     make(http.Header),
				}, nil
			},
		}
		slack := &Slack{
			Client: mockClient,
			Message: Message{
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
			},
		}
		errExpected := errors.New("error sending message: invalid_auth")

		err := slack.Send(context.Background())

		assert.AreEqualErrs(
			t,
			err,
			errExpected,
			"Expected error from Slack API",
		)
	})
}
