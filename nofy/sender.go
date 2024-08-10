package nofy

import (
	"context"
)

type Messenger interface {
	Send(ctx context.Context) error
}

type Sender struct {
	Messengers []Messenger
}

func NewSender(m ...Messenger) *Sender {
	return &Sender{
		Messengers: m,
	}
}

func (s *Sender) AddMessenger(m Messenger) {
	s.Messengers = append(s.Messengers, m)
}

func (s *Sender) RemoveMessenger(m Messenger) {
	for i, msgr := range s.Messengers {
		if msgr == m {
			s.Messengers = append(s.Messengers[:i], s.Messengers[i+1:]...)
			break
		}
	}
}

func (s *Sender) SendAll(ctx context.Context) error {
	for _, m := range s.Messengers {
		err := m.Send(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
