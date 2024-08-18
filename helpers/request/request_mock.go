package request

import (
	"context"
	"net/http"
)

type MockRequester struct {
	DoFunc func(ctx context.Context, options ...Option) (*http.Response, []byte, error)
}

func (m *MockRequester) Do(ctx context.Context, options ...Option) (*http.Response, []byte, error) {
	return m.DoFunc(ctx, options...)
}
