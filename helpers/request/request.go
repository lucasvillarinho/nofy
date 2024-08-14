package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Request struct {
	Headers map[string]string
	Payload map[string]any
	Client  *http.Client
	Timeout time.Duration
	Method  string
	URL     string
}

type Option func(*Request)

// WithMethod sets the method for the request.
func WithMethod(method string) Option {
	return func(r *Request) {
		r.Method = method
	}
}

// WithURL sets the URL for the request.
func WithURL(url string) Option {
	return func(r *Request) {
		r.URL = url
	}
}

// WithHeader sets a header for the request.
func WithHeader(key, value string) Option {
	return func(r *Request) {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		r.Headers[key] = value
	}
}

// WithTimeout sets the timeout for the request.
func WithTimeout(timeout time.Duration) Option {
	return func(r *Request) {
		r.Timeout = timeout
	}
}

// WithClient sets the client to use for the request.
func WithClient(client *http.Client) Option {
	return func(r *Request) {
		r.Client = client
	}
}

// WithPayload sets the payload of the request.
func WithPayload(payload map[string]any) Option {
	return func(r *Request) {
		r.Payload = payload
	}
}

// Do sends a request to the given URL with the given method, headers, and payload.
// It returns the response from the server.
// If the request fails, it returns an error.
func Do(options ...Option) (*http.Response, error) {
	r := &Request{}

	for _, opt := range options {
		opt(r)
	}

	if err := validate(r); err != nil {
		return nil, err
	}

	body, err := json.Marshal(r.Payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for header, headerValue := range r.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	return resp, nil
}

// any is a type that can hold any value.
func validate(r *Request) error {
	if r.Method == "" {
		return fmt.Errorf("method is required")
	}

	if r.URL == "" {
		return fmt.Errorf("url is required")
	}

	if r.Timeout == 0 {
		return fmt.Errorf("timeout is required")
	}

	if r.Client == nil {
		r.Client = &http.Client{
			Timeout: r.Timeout,
		}
	}

	return nil
}
