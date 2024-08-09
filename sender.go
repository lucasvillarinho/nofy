package nofy

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context) error
	AddRecipient(recipient any) error
	RemoveRecipient(recipient any) error
	GetId() string
}
