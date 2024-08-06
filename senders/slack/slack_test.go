package slack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlackClient(t *testing.T) {
	t.Run("Error when token is empty", func(t *testing.T) {
		slack, err := NewSlackClient("")

		require.Error(t, err)
		assert.Nil(t, slack)
		assert.Equal(t, "missing Slack token", err.Error())
	})

	t.Run("Create Slack client with default values", func(t *testing.T) {
		slack, err := NewSlackClient("mock-token")

		require.NoError(t, err)
		assert.Equal(t, "https://slack.com/api/chat.postMessage", slack.url)
		assert.Equal(t, "mock-token", slack.token)
		assert.Equal(t, 5000*time.Millisecond, slack.timeout)
		assert.Empty(t, slack.recipients)
	})

	t.Run("Create Slack client with custom timeout", func(t *testing.T) {
		slack, err := NewSlackClient("mock-token", WithTimeout(10*time.Second))
		slack.token = "mock-token"

		require.NoError(t, err)
		assert.Equal(t, 10*time.Second, slack.timeout)
	})
}

func TestSlacksend(t *testing.T) {
	t.Run("Successful message send", func(t *testing.T) {
		server := mockSlackServer(t, http.StatusOK, Response{OK: true})
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
		}

		err := slack.send(
			context.Background(),
			[]byte(`{"text":"Hello, world!"}`),
		)
		assert.NoError(t, err)
	})

	t.Run("Error creating request", func(t *testing.T) {
		slack := &Slack{
			url:    "://bad-url",
			token:  "test-token",
			client: &http.Client{},
		}

		err := slack.send(
			context.Background(),
			[]byte(`{"text":"Hello, world!"}`),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error creating request")
	})

	t.Run("Error sending message", func(t *testing.T) {
		server := mockSlackServer(
			t,
			http.StatusInternalServerError,
			Response{OK: false, Error: "internal_error"},
		)
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
		}

		err := slack.send(
			context.Background(),
			[]byte(`{"text":"Hello, world!"}`),
		)
		assert.Error(t, err)
		assert.Contains(
			t,
			err.Error(),
			"error sending message: 500 Internal Server Error",
		)
	})

	t.Run("Error reading response", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{invalid json}"))
			}),
		)
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
		}

		err := slack.send(
			context.Background(),
			[]byte(`{"text":"Hello, world!"}`),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshalling response")
	})

	t.Run("Error reading response body", func(t *testing.T) {
		server := mockSlackServerWithErrorOnRead(t)
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
		}

		err := slack.send(
			context.Background(),
			[]byte(`{"text":"Hello, world!"}`),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error reading response")
	})

	t.Run("Error in Slack response", func(t *testing.T) {
		server := mockSlackServer(
			t,
			http.StatusOK,
			Response{OK: false, Error: "invalid_auth"},
		)
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
		}

		err := slack.send(
			context.Background(),
			[]byte(`{"text":"Hello, world!"}`),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error sending message: invalid_auth")
	})
}

func TestRemoveRecipient(t *testing.T) {
	t.Run("Existing recipient", func(t *testing.T) {
		s := &Slack{
			recipients: []Recipient{
				{Channel: "channel1"},
				{Channel: "channel2"},
				{Channel: "channel3"},
			},
		}
		expectedRecipients := []Recipient{
			{Channel: "channel1"},
			{Channel: "channel3"},
		}

		s.RemoveRecipient("channel2")
		assert.Equal(t, expectedRecipients, s.recipients)
	})

	t.Run("Non-existing recipient", func(t *testing.T) {
		s := &Slack{
			recipients: []Recipient{
				{Channel: "channel1"},
				{Channel: "channel2"},
				{Channel: "channel3"},
			},
		}
		expectedRecipients := []Recipient{
			{Channel: "channel1"},
			{Channel: "channel2"},
			{Channel: "channel3"},
		}

		s.RemoveRecipient("channel4")
		assert.Equal(t, expectedRecipients, s.recipients)
	})
}

func TestAddRecipient(t *testing.T) {
	t.Run("Single recipient", func(t *testing.T) {
		s := &Slack{}
		expectedRecipients := []Recipient{
			{Channel: "channel1"},
		}

		s.AddRecipient("channel1")
		assert.Equal(t, expectedRecipients, s.recipients)
	})

	t.Run("Multiple recipients", func(t *testing.T) {
		s := &Slack{}
		expectedRecipients := []Recipient{
			{Channel: "channel1"},
			{Channel: "channel2"},
			{Channel: "channel3"},
		}

		s.AddRecipient("channel1")
		s.AddRecipient("channel2")
		s.AddRecipient("channel3")
		assert.Equal(t, expectedRecipients, s.recipients)
	})
}

func TestSendBlocks(t *testing.T) {
	t.Run("Successful send blocks to multiple recipients", func(t *testing.T) {
		server := mockSlackServer(t, http.StatusOK, Response{OK: true})
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
			recipients: []Recipient{
				{Channel: "channel1"},
				{Channel: "channel2"},
			},
		}

		blocks := []map[string]interface{}{
			{"type": "section", "text": map[string]string{"type": "mrkdwn", "text": "Hello, world!"}},
		}

		err := slack.SendBlocks(context.Background(), blocks)
		require.NoError(t, err)
	})

	t.Run("Error marshalling message", func(t *testing.T) {
		slack := &Slack{
			recipients: []Recipient{
				{Channel: "channel1"},
			},
		}

		blocks := []map[string]interface{}{
			{"type": "section", "text": make(chan int)}, // Invalid data type for JSON marshalling
		}

		err := slack.SendBlocks(context.Background(), blocks)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error marshalling message")
	})

	t.Run("Error sending message with blocks", func(t *testing.T) {
		server := mockSlackServer(t, http.StatusInternalServerError, Response{OK: false, Error: "internal_error"})
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
			recipients: []Recipient{
				{Channel: "channel1"},
			},
		}

		blocks := []map[string]interface{}{
			{"type": "section", "text": map[string]string{"type": "mrkdwn", "text": "Hello, world!"}},
		}

		err := slack.SendBlocks(context.Background(), blocks)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error sending message with blocks")
	})
}

func TestSlackSend(t *testing.T) {
	t.Run("Successful send message", func(t *testing.T) {
		server := mockSlackServer(t, http.StatusOK, Response{OK: true})
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
			recipients: []Recipient{
				{Channel: "channel1"},
			},
		}

		err := slack.Send(context.Background(), "Hello, world!")
		require.NoError(t, err)
	})

	t.Run("Error sending message with blocks", func(t *testing.T) {
		server := mockSlackServer(t, http.StatusInternalServerError, Response{OK: false, Error: "internal_error"})
		defer server.Close()

		slack := &Slack{
			url:    server.URL + "/api/chat.postMessage",
			token:  "test-token",
			client: server.Client(),
			recipients: []Recipient{
				{Channel: "channel1"},
			},
		}

		err := slack.Send(context.Background(), "Hello, world!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error sending message with blocks")
	})
}

type mockHTTPClient struct{}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("mock client error")
}

func mockSlackServer(
	t *testing.T,
	statusCode int,
	response interface{},
) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/chat.postMessage", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	})

	return httptest.NewServer(handler)
}

func mockSlackServerWithErrorOnRead(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/chat.postMessage", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")

		w.WriteHeader(http.StatusOK)
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Fatalf("could not hijack connection: %v", err)
		}
		conn.Close()
	})

	return httptest.NewServer(handler)
}
