package request

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lucasvillarinho/nofy/helpers/assert"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestValidate(t *testing.T) {
	t.Run("missing method", func(t *testing.T) {
		r := &Request{
			URL: "https://example.com",
		}

		err := validate(r)
		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "method is required")
	})
	t.Run("missing URL", func(t *testing.T) {
		r := &Request{
			Method: "GET",
		}

		err := validate(r)
		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "url is required")
	})
	t.Run("Valid", func(t *testing.T) {
		r := &Request{
			Method: "GET",
			URL:    "https://example.com"}

		err := validate(r)
		assert.IsNil(t, err)
	})
}

func TestWith(t *testing.T) {
	t.Run("with client", func(t *testing.T) {
		r := &Request{}
		client := &http.Client{}

		WithClient(client)(r)

		assert.AreEqual(t, r.Client, client, "Expected client to be set")
	})

	t.Run("with method", func(t *testing.T) {
		r := &Request{}

		WithMethod("GET")(r)

		assert.AreEqual(t, r.Method, "GET", "Expected method to be set")
	})

	t.Run("with URL", func(t *testing.T) {
		r := &Request{}

		WithURL("https://example.com")(r)

		assert.AreEqual(t, r.URL, "https://example.com", "Expected URL to be set")
	})

	t.Run("with header", func(t *testing.T) {
		r := &Request{}

		WithHeader("key", "value")(r)

		assert.AreEqual(t, r.Headers["key"], "value", "Expected header to be set")
	})

	t.Run("with payload", func(t *testing.T) {
		r := &Request{}
		payload := map[string]interface{}{
			"key": "value",
		}
		payloadBytes, _ := json.Marshal(payload)

		WithPayload(payloadBytes)(r)
	})
}

func TestDo(t *testing.T) {
	t.Run("invalid request", func(t *testing.T) {
		resp, err := DoWithContext(context.Background())

		assert.IsNil(t, resp)
		assert.AreEqualErrs(t, err, errors.New("method is required"), "Expected error")
	})

	t.Run("error creating request", func(t *testing.T) {
		resp, err := DoWithContext(
			context.TODO(),
			WithMethod(http.MethodGet),
			WithURL("://invalid-url"),
			WithClient(http.DefaultClient),
		)
		expectedErr := errors.New("error creating request: parse \"://invalid-url\": missing protocol scheme")

		assert.IsNil(t, resp)
		assert.AreEqualErrs(t, err, expectedErr, "Expected error creating request")
	})

	t.Run("error sending request", func(t *testing.T) {
		mockClient := MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}
		resp, err := DoWithContext(
			context.TODO(),
			WithMethod(http.MethodGet),
			WithURL("https://example.com"),
			WithClient(mockClient),
		)
		expectedErr := errors.New("error sending request: network error")

		assert.IsNil(t, resp)
		assert.AreEqualErrs(t, err, expectedErr, "Expected error sending request")
	})

	t.Run("success", func(t *testing.T) {
		mockClient := MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"message": "success"}`)),
					Header: map[string][]string{
						"Content-Type": {"application/json"},
					},
				}, nil
			},
		}
		resp, err := DoWithContext(
			context.Background(),
			WithMethod(http.MethodGet),
			WithURL("https://example.com"),
			WithClient(mockClient),
		)

		assert.AreEqual(t, resp.StatusCode, http.StatusOK, "Expected status code")
		body, _ := io.ReadAll(resp.Body)
		assert.AreEqual(t, `{"message": "success"}`, string(body))
		assert.AreEqual(t, resp.Header.Get("Content-Type"), "application/json", "Expected content type")
		assert.IsNil(t, err)
	})
}
