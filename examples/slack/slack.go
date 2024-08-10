package main

import (
	"context"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/messengers/slack"
)

func main() {
	// Create a new Slack messenger
	slackMensseger, _ := slack.NewSlackMensseger(
		// Set the Slack token to be used to send (required)
		slack.WithToken("test-token"),
		slack.WithMessage(
			// Message to be sent to the slack channel (required)
			// The message is a slice of maps, each map represents a block of the message
			// In this case, we are sending a single block with a text section
			slack.Message{
				Channel: "test-channel",
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
	nofy := nofy.NewWithMessengers(slackMensseger)

	// Send the message for all messengers
	err := nofy.SendAll(context.Background())
	if err != nil {
		panic(err)
	}
}
