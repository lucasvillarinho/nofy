package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/helpers/assert"
	"github.com/lucasvillarinho/nofy/messengers/slack"
)

func TestSend(t *testing.T) {
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	t.Run("should send message to slack", func(t *testing.T) {
		slackMessenger, _ := slack.NewSlackMessenger(
			slack.WithToken(slackToken),
			slack.WithMessage(
				slack.Message{
					Channel: slackChannel,
					Content: []map[string]interface{}{
						{
							"type": "section",
							"text": map[string]string{
								"type": "mrkdwn",
								"text": "Test Nofy",
							},
						},
					},
				}))
		nofy := nofy.NewWithMessengers(slackMessenger)

		err := nofy.SendAll(context.Background())

		assert.IsNil(t, err)
	})

}
