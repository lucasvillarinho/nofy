package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type request struct {
	headers map[string]string
	payload []byte
	client  HTTPClient
	method  string
	url     string
}

type Requester interface {
	Do(ctx context.Context, options ...Option) (*http.Response, []byte, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Option func(*request)

func NewRequester() Requester {
	return &request{
		client: http.DefaultClient,
	}
}

// WithMethod sets the method for the request.
func WithMethod(method string) Option {
	return func(r *request) {
		r.method = method
	}
}

// WithURL sets the URL for the request.
func WithURL(url string) Option {
	return func(r *request) {
		r.url = url
	}
}

// WithHeader sets a header for the request.
func WithHeader(key, value string) Option {
	return func(r *request) {
		if r.headers == nil {
			r.headers = make(map[string]string)
		}
		r.headers[key] = value
	}
}

// WithClient sets the client to use for the request.
func WithClient(client HTTPClient) Option {
	return func(r *request) {
		r.client = client
	}
}

// WithPayload sets the payload of the request.
func WithPayload(payload []byte) Option {
	return func(r *request) {
		r.payload = payload
	}
}

// Do sends a request to the given URL with the given method, headers, and payload.
// It returns the response from the server and the body of the response.
func (r *request) Do(ctx context.Context, options ...Option) (*http.Response, []byte, error) {
	rq := &request{}
	for _, opt := range options {
		opt(rq)
	}

	if err := validate(rq); err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		rq.method,
		rq.url,
		bytes.NewBuffer(rq.payload),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	for header, headerValue := range rq.headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := rq.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error sending request: %w", err)
	}

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}
	defer resp.Body.Close()

	return resp, bodyResponse, nil
}

// any is a type that can hold any value.
func validate(r *request) error {
	if r.method == "" {
		return fmt.Errorf("method is required")
	}

	if r.url == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}
