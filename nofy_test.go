package nofy

import (
	"context"
	"errors"
	"testing"

	"github.com/lucasvillarinho/nofy/helpers/assert"
)

type MockMessenger struct {
	sendFunc func(ctx context.Context) error
}

func (m *MockMessenger) Send(ctx context.Context) error {
	if m.sendFunc != nil {
		return m.sendFunc(ctx)
	}
	return nil
}

func TestNew(t *testing.T) {
	t.Run("should create a new Nofy instance with no messengers", func(t *testing.T) {
		nofy := New()

		assert.IsNotNil(t, nofy, "Expected Nofy instance to be created")
		assert.AreEqual(t, len(nofy.messengers), 0, "Expected messengers list to be empty")
	})
}

func TestNewWithMessengers(t *testing.T) {
	t.Run("should create a new Nofy instance with provided messengers", func(t *testing.T) {
		mockMessenger1 := &MockMessenger{}
		mockMessenger2 := &MockMessenger{}

		nofy := NewWithMessengers(mockMessenger1, mockMessenger2)

		assert.IsNotNil(t, nofy, "Expected Nofy instance to be created")
		assert.AreEqual(t, len(nofy.messengers), 2, "Expected messengers list to contain 2 items")
	})

	t.Run(
		"should create a new Nofy instance with an empty list when no messengers are provided",
		func(t *testing.T) {
			nofy := NewWithMessengers()

			assert.IsNotNil(t, nofy, "Expected Nofy instance to be created")
			assert.AreEqual(t, len(nofy.messengers), 0, "Expected messengers list to be empty")
		},
	)
}

func TestAddMessenger(t *testing.T) {
	t.Run("should add a messenger to the list", func(t *testing.T) {
		s := &Nofy{}
		m1 := &MockMessenger{}

		s.AddMessenger(m1)

		assert.AreEqual(
			t,
			s.messengers[0],
			m1,
			"Expected messenger to be added",
		)
	})

	t.Run("should add multiple messengers to the list", func(t *testing.T) {
		s := &Nofy{}
		m1 := &MockMessenger{}
		m2 := &MockMessenger{}

		s.AddMessenger(m1)
		s.AddMessenger(m2)

		assert.AreEqual(
			t,
			s.messengers[0],
			m1,
			"Expected first messenger to be added",
		)
		assert.AreEqual(
			t,
			s.messengers[1],
			m2,
			"Expected second messenger to be added",
		)
	})
}

func TestRemoveMessenger(t *testing.T) {
	t.Run("should remove messenger from middle of list", func(t *testing.T) {
		s := Nofy{}
		m1 := &MockMessenger{}
		m2 := &MockMessenger{}
		m3 := &MockMessenger{}
		s.AddMessenger(m1)
		s.AddMessenger(m2)
		s.AddMessenger(m3)

		s.RemoveMessenger(m2)

		assert.AreEqual(
			t,
			len(s.messengers),
			2,
			"Expected one messenger to be removed",
		)
		assert.AreEqual(
			t,
			s.messengers[0],
			m1,
			"Expected first messenger to remain",
		)
		assert.AreEqual(
			t,
			s.messengers[1],
			m3,
			"Expected last messenger to remain",
		)
	})

	t.Run("should not remove messenger not in list", func(t *testing.T) {
		s := Nofy{}
		m1 := &MockMessenger{}
		m2 := &MockMessenger{}
		m3 := &MockMessenger{}
		s.AddMessenger(m1)
		s.AddMessenger(m3)

		s.RemoveMessenger(m2)

		assert.AreEqual(
			t,
			len(s.messengers),
			2,
			"Expected no messenger to be removed",
		)
		assert.AreEqual(
			t,
			s.messengers[0],
			m1,
			"Expected first messenger to remain",
		)
		assert.AreEqual(
			t,
			s.messengers[1],
			m3,
			"Expected last messenger to remain",
		)
	})

	t.Run("should remove first messenger from list", func(t *testing.T) {
		s := Nofy{}
		m1 := &MockMessenger{}
		m2 := &MockMessenger{}
		s.AddMessenger(m1)
		s.AddMessenger(m2)

		s.RemoveMessenger(m1)

		assert.AreEqual(
			t,
			len(s.messengers),
			1,
			"Expected one messenger to be removed",
		)
		assert.AreEqual(
			t,
			s.messengers[0],
			m2,
			"Expected last messenger to remain",
		)
	})

	t.Run("should remove last messenger from list", func(t *testing.T) {
		s := Nofy{}
		m1 := &MockMessenger{}
		s.AddMessenger(m1)

		s.RemoveMessenger(m1)

		assert.AreEqual(
			t,
			len(s.messengers),
			0,
			"Expected one messenger to be removed",
		)
	})
}

func TestAggregateErrors(t *testing.T) {
	t.Run("should return nil when no errors are present", func(t *testing.T) {
		errChan := make(chan error)
		close(errChan)

		err := aggregateErrors(errChan)

		assert.IsNil(t, err, "Expected no errors")
	})

	t.Run("should return an error when one error is present", func(t *testing.T) {
		errChan := make(chan error, 1)
		errChan <- errors.New("single error occurred")
		close(errChan)
		expectedErr := errors.New("errors: single error occurred")

		err := aggregateErrors(errChan)

		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected single error",
		)
	})

	t.Run("should return an error when multiple errors are present", func(t *testing.T) {
		errChan := make(chan error, 2)
		errChan <- errors.New("first error")
		errChan <- errors.New("second error")
		close(errChan)
		expectedErr := errors.New("errors: first error; second error")

		err := aggregateErrors(errChan)

		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected multiple errors",
		)
	})
}

func TestSendAll(t *testing.T) {
	t.Run("should return nil when messages are sent successfully", func(t *testing.T) {
		s := &Nofy{
			messengers: []Messenger{
				&MockMessenger{},
				&MockMessenger{},
			},
		}

		err := s.SendAll(context.Background())

		assert.IsNil(t, err, "Expected no errors")
	})

	t.Run("should return an error when one messenger fails", func(t *testing.T) {
		s := &Nofy{
			messengers: []Messenger{
				&MockMessenger{},
				&MockMessenger{
					sendFunc: func(ctx context.Context) error {
						return errors.New("failed to send message")
					},
				},
			},
		}
		expectedErr := errors.New("errors: failed to send message")

		err := s.SendAll(context.Background())

		assert.AreEqualErrs(
			t,
			err,
			expectedErr,
			"Expected one error",
		)
	})

	t.Run("should return an error error when all messages fail", func(t *testing.T) {
		s := &Nofy{
			messengers: []Messenger{
				&MockMessenger{
					sendFunc: func(_ context.Context) error {
						return errors.New("first message failed")
					},
				},
				&MockMessenger{
					sendFunc: func(_ context.Context) error {
						return errors.New("second message failed")
					},
				},
			},
		}

		err := s.SendAll(context.Background())

		assert.IsNotNil(t, err, "Expected errors")
	})

	t.Run("should handle panic gracefully and return error", func(t *testing.T) {
		s := &Nofy{
			messengers: []Messenger{
				&MockMessenger{
					sendFunc: func(_ context.Context) error {
						panic("unexpected panic")
					},
				},
				&MockMessenger{},
			},
		}
		expectedErr := errors.New("errors: panic recovered: unexpected panic")

		err := s.SendAll(context.Background())

		assert.IsNotNil(t, err, "Expected an error due to panic")
		assert.AreEqualErrs(t, err, expectedErr, "Expected panic to be recovered")
	})
}
