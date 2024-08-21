package main

import (
	"context"

	"github.com/lcvilla/nofy"
	"github.com/lcvilla/nofy/messengers/slack"
)

func main() {
	// Create a new Slack messenger
	slackMessenger, _ := slack.NewSlackMessenger(
		// Set the Slack token to be used to send (required)
		slack.WithToken("token"),
		slack.WithMessage(
			// Message to be sent to the slack channel (required)
			// The message is a slice of maps, each map represents a block of the message
			// In this case, we are sending a single block with a text section
			slack.Message{
				Channel: "channelID",
				Content: []map[string]interface{}{
					{
						"type": "section",
						"text": map[string]string{
							"type": "mrkdwn",
							"text": "Hello, World!",
						},
					},
				},
			}))

	// Create a new Nofy with the Slack messenger
	nofy := nofy.NewWithMessengers(slackMessenger)

	// Send the message for all messengers
	err := nofy.SendAll(context.Background())
	if err != nil {
		panic(err)
	}
}
