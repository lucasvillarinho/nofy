package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ErrorRoundTripper struct{}

func (c *ErrorRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	return nil, errors.New("simulated network error")
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewSlackClient(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		wantErr error
	}{
		{
			name: "missing token",
			options: []Option{
				WithRecipients([]Recipient{{}}),
				WithMessage([]map[string]any{{"text": "Test message"}}),
			},
			wantErr: fmt.Errorf("missing Slack Token"),
		},
		{
			name: "missing recipients",
			options: []Option{
				WithToken("test-token"),
				WithMessage([]map[string]any{{"text": "Test message"}}),
			},
			wantErr: fmt.Errorf("missing Slack Recipients"),
		},
		{
			name: "missing message",
			options: []Option{
				WithToken("test-token"),
				WithRecipients([]Recipient{{}}),
			},
			wantErr: fmt.Errorf("missing Slack message"),
		},
		{
			name: "missing timeout",
			options: []Option{
				WithTimeout(0),
				WithToken(
					"test-token",
				), WithRecipients([]Recipient{{Channel: "test-channel"}}), WithMessage([]map[string]any{{"text": "Test message"}}),
			},
			wantErr: fmt.Errorf("missing Timeout"),
		},
		{
			name: "all options provided",
			options: []Option{
				WithToken("test-token"),
				WithRecipients(
					[]Recipient{{}},
				), WithMessage([]map[string]any{{"text": "Test message"}}),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSlackClient(tt.options...)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSlackSend(t *testing.T) {
	t.Run("error creating request", func(t *testing.T) {
		client := &http.Client{}
		slack := &Slack{
			URL:    "://invalid-url",
			Token:  "test-token",
			Client: client,
		}
		body := []byte(`{"text":"Hello, World!"}`)

		resp, err := slack.send(context.Background(), body)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "error creating request")
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"ok": false}`))
			}),
		)
		defer server.Close()

		slack := &Slack{
			URL:   server.URL,
			Token: "test-token",
		}

		body := []byte(`{"text":"Hello, World!"}`)
		resp, err := slack.send(context.Background(), body)

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(
			t,
			"error sending message: 500 Internal Server Error",
			err.Error(),
		)
	})

	t.Run("client do error", func(t *testing.T) {
		client := &http.Client{
			Transport: &ErrorRoundTripper{},
		}
		slack := &Slack{
			URL:    "http://example.com",
			Token:  "test-token",
			Client: client,
		}
		body := []byte(`{"text":"Hello, World!"}`)

		resp, err := slack.send(context.Background(), body)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(
			t,
			"error sending message: Post \"http://example.com\": simulated network error",
			err.Error(),
		)
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`invalid json`))
			}),
		)
		defer server.Close()

		slack := &Slack{
			URL:   server.URL,
			Token: "test-token",
		}
		body := []byte(`{"text":"Hello, World!"}`)

		resp, err := slack.send(context.Background(), body)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(
			t,
			"error unmarshalling response: invalid character 'i' looking for beginning of value",
			err.Error(),
		)
	})

	t.Run("successful message send", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"ok": true}`))
			}),
		)
		defer server.Close()

		slack := &Slack{
			URL:   server.URL,
			Token: "test-token",
		}
		body := []byte(`{"text":"Hello, World!"}`)

		resp, err := slack.send(context.Background(), body)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.OK)
	})
}

func TestSend(t *testing.T) {
	t.Run("Successful message send", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		s := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "valid-token",
			Client: mockClient,
			Recipients: []Recipient{
				{Channel: "channel1"},
			},
			Message: []map[string]any{
				{"type": "section", "text": "Hello, world!"},
			},
		}

		slackResponse := Response{
			OK: true,
		}
		responseBody, _ := json.Marshal(slackResponse)
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
		}

		mockClient.
			On("Do", mock.Anything).
			Return(mockResponse, nil).Once()

		err := s.Send(context.Background())
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error marshalling message", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		s := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "valid-token",
			Client: mockClient,
			Recipients: []Recipient{
				{Channel: "channel1"},
			},
			Message: []map[string]any{
				{"type": "section", "text": make(chan int)},
			},
		}

		err := s.Send(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error marshalling message")
	})

	t.Run("Slack response with error", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		s := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "valid-token",
			Client: mockClient,
			Recipients: []Recipient{
				{Channel: "channel1"},
			},
			Message: []map[string]any{
				{"type": "section", "text": "Hello, world!"},
			},
		}

		slackResponse := Response{
			OK:    false,
			Error: "invalid_auth",
		}
		responseBody, _ := json.Marshal(slackResponse)
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
		}

		mockClient.
			On("Do", mock.Anything).
			Return(mockResponse, nil).Once()

		err := s.Send(context.Background())
		assert.Error(t, err)
		assert.Contains(
			t,
			err.Error(),
			"error sending message with blocks: invalid_auth",
		)
	})

	t.Run("Error waiting for group", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mockClient := new(MockHTTPClient)
		s := &Slack{
			URL:    "https://slack.com/api/chat.postMessage",
			Token:  "valid-token",
			Client: mockClient,
			Recipients: []Recipient{
				{Channel: "channel1"},
			},
			Message: []map[string]any{
				{"type": "section", "text": "Hello, world!"},
			},
		}

		err := s.Send(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error sending message with blocks")
	})
}
