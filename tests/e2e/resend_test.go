package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/helpers/assert"
	"github.com/lucasvillarinho/nofy/messengers/resend"
)

func TestResend(t *testing.T) {
	resendToken := os.Getenv("RESEND_TOKEN")
	resendTo := os.Getenv("RESEND_TO")

	if resendToken == "" || resendTo == "" {
		t.Fatal(
			"E2E Test Setup: Environment variables RESEND_TOKEN and RESEND_TO must be set before running the end-to-end tests.",
		)
	}

	t.Run("should send message to resend", func(t *testing.T) {
		resendMessenger, err := resend.NewResendMessenger(
			resend.WithToken(resendToken),
			resend.WithMessage(
				&resend.Message{
					From:    "onboarding@resend.dev",
					To:      []string{resendTo},
					Subject: "Test e2e - NoFy",
					HTML:    "<p> Test e2e</p>",
				}),
		)

		assert.IsNil(t, err)

		nofy := nofy.NewWithMessengers(resendMessenger)

		err = nofy.SendAll(context.Background())

		assert.IsNil(t, err)
	})
}
