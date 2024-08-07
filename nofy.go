package main

import (
	"context"

	"github.com/lucasvillarinho/nofy/models"
)

type Nofy struct {
	senders []models.Sender
}

func NewNofy() *Nofy {
	return &Nofy{
		senders: make([]models.Sender, 0),
	}
}

func (n *Nofy) AddSender(s models.Sender) {
	n.senders = append(n.senders, s)
}

func (n *Nofy) RemoveSender(s models.Sender) {
	for i, sender := range n.senders {
		if sender == s {
			n.senders = append(n.senders[:i], n.senders[i+1:]...)
			return
		}
	}
}

func (n *Nofy) Send(ctx context.Context, message any) {
	for _, sender := range n.senders {
		sender.Send(ctx)
	}
}
