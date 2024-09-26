package resend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/lucasvillarinho/nofy/helpers/assert"
	"github.com/lucasvillarinho/nofy/helpers/request"
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

func TestSend(t *testing.T) {
	t.Run("should send message successfully", func(t *testing.T) {
		mockRequester := &request.MockRequester{
			DoFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, []byte(`{"ok": true}`), nil
			},
		}
		message := Message{
			From:    "test-from",
			To:      []string{"test-to"},
			Subject: "test-subject",
			HTML:    "<p> Text Html</p>",
		}

		messenger := &Resend{
			Token:     "test-token",
			URL:       "https://api.sendgrid.com/v3/mail/send",
			Timeout:   5 * time.Second,
			Message:   message,
			requester: mockRequester,
		}

		err := messenger.Send(context.TODO())

		assert.IsNil(t, err)
	})

	t.Run("should return error when marshalling message fails", func(t *testing.T) {
		MarshalFunc = func(_ interface{}) ([]byte, error) {
			return nil, errors.New("invalid payload")
		}

		message := Message{
			From:    "test-from",
			To:      []string{},
			Subject: "\x80\x81",
			HTML:    fmt.Sprintf("%f", math.NaN()),
		}

		messenger := &Resend{
			Token:     "test-token",
			Timeout:   5 * time.Second,
			URL:       "https://api.resend.com/emails",
			Message:   message,
			requester: nil,
		}

		err := messenger.Send(context.TODO())

		assert.AreEqualErrs(
			t,
			err,
			errors.New("error marshaling message: invalid payload"),
			"Expected error marshaling message",
		)

		MarshalFunc = func(v interface{}) ([]byte, error) {
			return json.Marshal(v)
		}
	})

	t.Run("should return error when request fails", func(t *testing.T) {
		mockRequester := &request.MockRequester{
			DoFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
				return nil, nil, errors.New("error sending request")
			},
		}
		message := Message{
			From:    "test-from",
			To:      []string{"test-to"},
			Subject: "test-subject",
			HTML:    "<p> Text Html</p>",
		}

		messenger := &Resend{
			Token:     "test-token",
			URL:       "https://api.resend.com/emails",
			Timeout:   5 * time.Second,
			Message:   message,
			requester: mockRequester,
		}

		err := messenger.Send(context.TODO())

		assert.AreEqualErrs(
			t,
			err,
			errors.New("error sending message: error sending request"),
			"Expected error sending request",
		)
	})

	t.Run("should return error when response is not OK", func(t *testing.T) {
		bodyResponse := `{"name": "missing_required_field","statusCode": 422,"message": "Missing from field."}`
		mockRequester := &request.MockRequester{
			DoFunc: func(ctx context.Context, options ...request.Option) (*http.Response, []byte, error) {
				return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
					}, []byte(bodyResponse),
					nil
			}}

		message := Message{
			From:    "test-from",
			To:      []string{"test-to"},
			Subject: "test-subject",
			HTML:    "<p> Text Html</p>",
		}

		messenger := &Resend{
			Token:     "test-token",
			URL:       "https://api.resend.com/emails",
			Timeout:   5 * time.Second,
			Message:   message,
			requester: mockRequester,
		}

		err := messenger.Send(context.TODO())

		assert.AreEqualErrs(
			t,
			err,
			fmt.Errorf("error sending message: status-code: %d body: %s",
				http.StatusUnprocessableEntity, bodyResponse),
			"Expected error sending message",
		)
	})

}
