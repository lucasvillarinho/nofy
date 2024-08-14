package request

import (
	"net/http"
	"testing"
	"time"

	"github.com/lucasvillarinho/nofy/helpers/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Valid request", func(t *testing.T) {
		r := &Request{
			Method:  "GET",
			URL:     "https://example.com",
			Timeout: 500 * time.Second,
		}

		err := validate(r)
		assert.IsNil(t, err)
	})
	t.Run("Invalid request - missing method", func(t *testing.T) {
		r := &Request{
			URL:     "https://example.com",
			Timeout: 500 * time.Second,
		}

		err := validate(r)
		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "method is required")
	})
	t.Run("Invalid request - missing URL", func(t *testing.T) {
		r := &Request{
			Method:  "GET",
			Timeout: 500 * time.Second,
		}

		err := validate(r)
		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "url is required")
	})
	t.Run("Invalid request - missing timeout", func(t *testing.T) {
		r := &Request{
			Method: "GET",
			URL:    "https://example.com",
		}

		err := validate(r)
		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "timeout is required")
	})
}

func TestWith(t *testing.T) {
	t.Run("With client", func(t *testing.T) {
		r := &Request{}
		client := &http.Client{}

		WithClient(client)(r)

		assert.AreEqual(t, r.Client, client, "Expected client to be set")
	})

	t.Run("With method", func(t *testing.T) {
		r := &Request{}

		WithMethod("GET")(r)

		assert.AreEqual(t, r.Method, "GET", "Expected method to be set")
	})

	t.Run("With URL", func(t *testing.T) {
		r := &Request{}

		WithURL("https://example.com")(r)

		assert.AreEqual(t, r.URL, "https://example.com", "Expected URL to be set")
	})

	t.Run("With header", func(t *testing.T) {
		r := &Request{}

		WithHeader("key", "value")(r)

		assert.AreEqual(t, r.Headers["key"], "value", "Expected header to be set")
	})

	t.Run("With timeout", func(t *testing.T) {
		r := &Request{}

		WithTimeout(500 * time.Second)(r)

		assert.AreEqual(t, r.Timeout, 500*time.Second, "Expected timeout to be set")
	})

	t.Run("With payload", func(t *testing.T) {
		r := &Request{}
		payload := map[string]interface{}{
			"key": "value",
		}

		WithPayload(payload)(r)

		assert.AreEqual(t, r.Payload["key"], "value", "Expected payload to be set")
	})
}
