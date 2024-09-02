package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("forced read error")
}

func TestValidate(t *testing.T) {
	t.Run("should pass validation when method and URL are present", func(t *testing.T) {
		r := &request{
			method: "GET",
			url:    "https://example.com",
		}

		err := validate(r)

		assert.IsNil(t, err)
	})
	t.Run("should return error when method is missing", func(t *testing.T) {
		r := &request{
			url: "https://example.com",
		}

		err := validate(r)

		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "method is required")
	})
	t.Run("should return error when URL is missing", func(t *testing.T) {
		r := &request{
			method: "GET",
		}

		err := validate(r)

		assert.IsNotNil(t, err)
		assert.AreEqual(t, err.Error(), "url is required")
	})
}

func TestWith(t *testing.T) {
	t.Run("should set client with WithClient option", func(t *testing.T) {
		r := &request{}
		client := &http.Client{}

		WithClient(client)(r)

		assert.AreEqual(t, r.client, client, "Expected client to be set")
	})

	t.Run("should set method with WithMethod option", func(t *testing.T) {
		r := &request{}

		WithMethod("GET")(r)

		assert.AreEqual(t, r.method, "GET", "Expected method to be set")
	})

	t.Run("should set URL with WithURL option", func(t *testing.T) {
		r := &request{}

		WithURL("https://example.com")(r)

		assert.AreEqual(
			t,
			r.url,
			"https://example.com",
			"Expected URL to be set",
		)
	})

	t.Run("should set header with WithHeader option", func(t *testing.T) {
		r := &request{}

		WithHeader("key", "value")(r)

		assert.AreEqual(
			t,
			r.headers["key"],
			"value",
			"Expected header to be set",
		)
	})

	t.Run("should set payload with WithPayload option", func(t *testing.T) {
		r := &request{}
		payload := map[string]interface{}{
			"key": "value",
		}
		payloadBytes, _ := json.Marshal(payload)

		WithPayload(payloadBytes)(r)
	})
}

func TestDo(t *testing.T) {
	t.Run("should return response successfully", func(t *testing.T) {
		mockClient := MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader(`{"message": "success"}`),
					),
					Header: map[string][]string{
						"Content-Type": {"application/json"},
					},
				}, nil
			},
		}

		resp, body, err := NewRequester().Do(
			context.Background(),
			WithMethod(http.MethodGet),
			WithURL("https://example.com"),
			WithClient(mockClient),
			WithHeader("Content-Type", "application/json"),
		)

		assert.AreEqual(
			t,
			resp.StatusCode,
			http.StatusOK,
			"Expected status code",
		)
		assert.AreEqual(t, `{"message": "success"}`, string(body))
		assert.AreEqual(
			t,
			resp.Header.Get("Content-Type"),
			"application/json",
			"Expected content type",
		)
		assert.IsNil(t, err)
	})

	t.Run("should return error when method is missing", func(t *testing.T) {
		expectedErr := errors.New("method is required")

		resp, body, err := NewRequester().Do(context.Background())

		assert.IsNil(t, resp)
		assert.IsNil(t, body)
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error",
		)
	})

	t.Run("should return error when creating request fails", func(t *testing.T) {
		mockClient := MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}
		expectedErr := errors.New(
			"error creating request: parse \"://invalid-url\": missing protocol scheme",
		)

		resp, body, err := NewRequester().Do(
			context.TODO(),
			WithMethod(http.MethodGet),
			WithURL("://invalid-url"),
			WithClient(mockClient),
		)

		assert.IsNil(t, resp)
		assert.IsNil(t, body)
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error sending request",
		)
	})

	t.Run("should return error when sending request fails", func(t *testing.T) {
		mockClient := MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}
		expectedErr := errors.New("error sending request: network error")

		resp, body, err := NewRequester().Do(
			context.TODO(),
			WithMethod(http.MethodGet),
			WithURL("https://example.com"),
			WithClient(mockClient),
		)

		assert.IsNil(t, resp)
		assert.IsNil(t, body)
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error sending request",
		)
	})

	t.Run("should return error when reading response body fails", func(t *testing.T) {
		mockClient := MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(&errorReader{}),
				}, nil
			},
		}
		expectedErr := errors.New("error reading response body: forced read error")

		resp, body, err := NewRequester().Do(
			context.Background(),
			WithMethod(http.MethodGet),
			WithURL("https://example.com"),
			WithClient(mockClient),
			WithHeader("Content-Type", "application/json"),
		)

		assert.IsNil(t, resp)
		assert.IsNil(t, body)
		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected error reading response body",
		)
	})
}
