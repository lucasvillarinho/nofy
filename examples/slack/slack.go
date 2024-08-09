package main

import (
	"context"

	"github.com/lucasvillarinho/nofy"
	"github.com/lucasvillarinho/nofy/senders/slack"
)

func main() {
	slackSender, err := slack.NewSlackClient(
		slack.WithToken("token"),
		slack.WithTimeout(5000))
	if err != nil {
		panic(err)
	}
	slackSender.AddRecipient("channel-1")

	nofy := nofy.NewNofy()
	nofy.AddSender(slackSender)
	err = nofy.Send(context.Background())
	if err != nil {
		panic(err)
	}
}
