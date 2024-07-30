package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackSendSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(Response{OK: true})
		w.Write(resp)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
	}

	err := s.send(context.Background(), []byte(`{"message":"test"}`))

	assert.NoError(t, err)
}

func TestSlackSendErrorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(Response{OK: false, Error: "some error"})
		w.Write(resp)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
	}

	err := s.send(context.Background(), []byte(`{"message":"test"}`))

	assert.Error(t, err)
	assert.EqualError(t, err, "error sending message: some error")
}

func TestSlackSendInvalidJSONResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
	}

	err := s.send(context.Background(), []byte(`{"message":"test"}`))

	assert.Error(t, err)
	assert.EqualError(t, err, "error unmarshalling response: invalid character 'i' looking for beginning of value")
}

func TestSlackSendHTTPRequestError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
	}

	err := s.send(context.Background(), []byte(`{"message":"test"}`))

	assert.Error(t, err)
	assert.EqualError(t, err, "error sending message: 500 Internal Server Error")
}

func TestSendBlocksSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(Response{OK: true})
		w.Write(resp)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
		recipients: []Recipient{
			{Channel: "test-channel"},
		},
	}

	blocks := []map[string]interface{}{
		{
			"type": "section",
			"text": map[string]string{
				"type": "mrkdwn",
				"text": "Test message",
			},
		},
	}

	err := s.SendBlocks(context.Background(), blocks)

	assert.NoError(t, err)
}

func TestSendBlocksMarshallingError(t *testing.T) {
	s := &Slack{
		url:   "http://invalid-url",
		token: "dummy-token",
		recipients: []Recipient{
			{Channel: "test-channel"},
		},
	}
	blocks := []map[string]interface{}{
		{
			"type": "section",
			"text": make(chan int),
		},
	}

	err := s.SendBlocks(context.Background(), blocks)

	assert.Error(t, err)
	assert.EqualError(t, err, "error marshalling message: json: unsupported type: chan int")
}

func TestSendBlocksSendError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
		recipients: []Recipient{
			{Channel: "test-channel"},
		},
	}
	blocks := []map[string]interface{}{
		{
			"type": "section",
			"text": map[string]string{
				"type": "mrkdwn",
				"text": "Test message",
			},
		},
	}

	err := s.SendBlocks(context.Background(), blocks)

	assert.Error(t, err)
	assert.EqualError(t, err, "error sending message with blocks: error sending message: 500 Internal Server Error")
}

func TestSendSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(Response{OK: true})
		w.Write(resp)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
		recipients: []Recipient{
			{Channel: "test-channel"},
		},
	}

	err := s.Send(context.Background(), "Test message")

	assert.NoError(t, err)
}

func TestSendSendError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	s := &Slack{
		url:   server.URL,
		token: "dummy-token",
		recipients: []Recipient{
			{Channel: "test-channel"},
		},
	}

	err := s.Send(context.Background(), "Test message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error sending message")
}

func TestRemoveRecipientExistingRecipient(t *testing.T) {
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
}

func TestRemoveRecipientNonExistingRecipient(t *testing.T) {
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
}

func TestAddRecipientSingleRecipient(t *testing.T) {
	s := &Slack{}
	expectedRecipients := []Recipient{
		{Channel: "channel1"},
	}

	s.AddRecipient("channel1")

	assert.Equal(t, expectedRecipients, s.recipients)
}

func TestAddRecipientMultipleRecipients(t *testing.T) {
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
}
