package request

import (
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
