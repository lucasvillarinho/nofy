package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type Requester interface {
	DoWithCtx(
		ctx context.Context,
		options ...Option,
	) (*http.Response, []byte, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Request struct {
	Headers map[string]string
	Payload []byte
	Client  HTTPClient
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

// WithClient sets the client to use for the request.
func WithClient(client HTTPClient) Option {
	return func(r *Request) {
		r.Client = client
	}
}

// WithPayload sets the payload of the request.
func WithPayload(payload []byte) Option {
	return func(r *Request) {
		r.Payload = payload
	}
}

// Do sends a request to the given URL with the given method, headers, and payload.
// It returns the response from the server.
// If the request fails, it returns an error.
func DoWithCtx(
	ctx context.Context,
	options ...Option,
) (*http.Response, []byte, error) {
	r := &Request{}

	for _, opt := range options {
		opt(r)
	}

	if err := validate(r); err != nil {
		return nil, nil, err
	}

	if r.Client == nil {
		r.Client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(
		ctx,
		r.Method,
		r.URL,
		bytes.NewBuffer(r.Payload),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	for header, headerValue := range r.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error sending request: %w", err)
	}

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response: %w", err)
	}
	defer resp.Body.Close()
	return resp, bodyResponse, nil
}

// any is a type that can hold any value.
func validate(r *Request) error {
	if r.Method == "" {
		return fmt.Errorf("method is required")
	}

	if r.URL == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}
