package nofy

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Messenger interface {
	Send(ctx context.Context) error
}

type Nofy struct {
	messengers []Messenger
}

func New() *Nofy {
	return &Nofy{
		messengers: make([]Messenger, 0),
	}
}

func NewWithMessengers(messengers ...Messenger) *Nofy {
	return &Nofy{
		messengers: messengers,
	}
}

func (s *Nofy) AddMessenger(m Messenger) {
	s.messengers = append(s.messengers, m)
}

func (s *Nofy) RemoveMessenger(m Messenger) {
	for i, msgr := range s.messengers {
		if msgr == m {
			s.messengers = append(s.messengers[:i], s.messengers[i+1:]...)
			break
		}
	}
}

func (s *Nofy) SendAll(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.messengers))

	for _, m := range s.messengers {
		wg.Add(1)
		go func(m Messenger) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic recovered: %v", r)
				}
			}()
			if err := m.Send(ctx); err != nil {
				errChan <- err
			}
		}(m)
	}

	wg.Wait()
	close(errChan)

	return aggregateErrors(errChan)
}

func aggregateErrors(errChan <-chan error) error {
	var errors []string
	for err := range errChan {
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(
			"errors: %s",
			strings.Join(errors, "; "),
		)
	}

	return nil
}
