package resend

import (
	"context"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/messengers/resend"
)

func main() {
	resendMessenger, _ := resend.NewResendMessenger(
		resend.WithToken("token"),
		resend.WithMessage(
			&resend.Message{
				From:    "test-from",
				To:      []string{"test-to"},
				Subject: "test-subject",
				HTML:    "<p> Text Html</p>",
			}),
	)

	nofy := nofy.NewWithMessengers(resendMessenger)

	err := nofy.SendAll(context.Background())
	if err != nil {
		panic(err)
	}
}
