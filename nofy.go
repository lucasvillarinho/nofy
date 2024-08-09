package nofy

import (
	"context"
	"fmt"

	pool "github.com/alitto/pond"
)

type Nofy struct {
	senders []Sender
}

func NewNofy() *Nofy {
	return &Nofy{
		senders: make([]Sender, 0),
	}
}

func (n *Nofy) AddSender(s Sender) {
	n.senders = append(n.senders, s)
}

func (n *Nofy) RemoveSender(s Sender) {
	for i, sender := range n.senders {
		if sender == s {
			n.senders = append(n.senders[:i], n.senders[i+1:]...)
			return
		}
	}
}

func (n *Nofy) Send(ctx context.Context) error {
	pool := pool.New(len(n.senders), len(n.senders))
	group, ctx := pool.GroupContext(ctx)

	for _, sender := range n.senders {
		group.Submit(func() error {
			err := sender.Send(ctx)
			if err != nil {
				return fmt.Errorf(
					"error sending message senderID: %s err: %v",
					sender.GetId(),
					err,
				)
			}

			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}
