package resend

import (
	"errors"
	"testing"
	"time"

	"github.com/lucasvillarinho/nofy/helpers/assert"
)

func TestNewResendMessenger(t *testing.T) {
	t.Run("should create Resend messenger successfully", func(t *testing.T) {
		messenger, err := NewResendMessenger(
			WithToken("test-token"),
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					From:    "test-from",
					To:      []string{"test-to"},
					Subject: "test-subject",
					HTML:    "<p> Text Html</p>",
				}),
		)

		assert.IsNil(t, err)
		assert.IsNotNil(t, messenger)
	})

	t.Run("should return error when token is missing", func(t *testing.T) {
		_, err := NewResendMessenger(
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					From:    "test-from",
					To:      []string{"test-to"},
					Subject: "test-subject",
					HTML:    "<p> Text Html</p>",
				}),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing token"),
			"Expected missing Resend Token error",
		)
	})

	t.Run("should return error when from is missing", func(t *testing.T) {
		_, err := NewResendMessenger(
			WithToken("test-token"),
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					To:      []string{"test-to"},
					Subject: "test-subject",
					HTML:    "<p> Text Html</p>",
				}),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing from"),
			"Expected missing Resend From error",
		)
	})

	t.Run("should return error when to is missing", func(t *testing.T) {
		_, err := NewResendMessenger(
			WithToken("test-token"),
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					From:    "test-from",
					Subject: "test-subject",
					HTML:    "<p> Text Html</p>",
				}),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing to"),
			"Expected missing Resend To error",
		)

	})

	t.Run("should return error when subject is missing", func(t *testing.T) {
		_, err := NewResendMessenger(
			WithToken("test-token"),
			WithTimeout(5*time.Second),
			WithMessage(
				Message{
					From: "test-from",
					To:   []string{"test-to"},
					HTML: "<p> Text Html</p>",
				}),
		)

		assert.AreEqualErrs(
			t,
			err,
			errors.New("missing subject"),
			"Expected missing Resend Subject error",
		)
	})
}
