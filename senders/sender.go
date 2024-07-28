package senders

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context, message string) error
	AddRecipients(channel string)
	RemoveRecipients(channel string)
}
