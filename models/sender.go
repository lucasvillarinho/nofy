package models

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context) error
}
