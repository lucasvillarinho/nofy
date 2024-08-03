package models

import (
	"context"
	"time"
)

type Sender interface {
	Send(ctx context.Context, message any) (*SendResult, error)
	AddRecipients(recipient any)
	RemoveRecipients(recipient any)
}

type SendResult struct {
	StatusCode   int
	ResponseTime time.Duration
	ResponseSize int64
	Error        error
}
