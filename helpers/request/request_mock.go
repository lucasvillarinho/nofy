package request

import (
	"context"
	"net/http"
)

type MockRequester struct {
	DoWithCtxFunc func(ctx context.Context, options ...Option) (*http.Response, []byte, error)
}

func (m *MockRequester) DoWithCtx(
	ctx context.Context,
	options ...Option,
) (*http.Response, []byte, error) {
	return m.DoWithCtxFunc(ctx, options...)
}
