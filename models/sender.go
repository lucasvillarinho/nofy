package models

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context) error
	AddMessage(message any) error
	AddRecipients(recipient any)
	RemoveRecipients(recipient any)
}
