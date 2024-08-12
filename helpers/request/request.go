package request

import (
	"fmt"
	"net/http"
	"time"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Client  *http.Client
	Timeout *time.Duration
}

type Option func(*Request)

func NewRequest(options ...Option) (*Request, error) {
	r := &Request{}

	for _, opt := range options {
		opt(r)
	}

	if r.Method == "" {
		return nil, fmt.Errorf("missing method")
	}

	if r.URL == "" {
		return nil, fmt.Errorf("missing URL")
	}

	if r.Timeout == nil {
		return nil, fmt.Errorf("missing timeout")
	}

	if r.Client == nil {
		r.Client = &http.Client{
			Timeout: *r.Timeout,
		}
	}

	return r, nil
}
