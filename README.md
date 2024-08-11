<p align="center"><img src="docs/images/logo.png" alt="nofylogo logo">
</h1>
<div align="center">
📬 NoFy is a versatile, zero-dependencies library for sending notifications to popular services.
<br>
<br>

![Zero Dependencies](https://img.shields.io/badge/Dependencies-Zero-brightgreen.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/lucasvillarinho/nofy)](https://goreportcard.com/report/github.com/lucasvillarinho/nofy) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/ec1e325348344d43906561ec19471598)](https://app.codacy.com/gh/lucasvillarinho/nofy/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)<a href="https://codeclimate.com/github/lucasvillarinho/nofy/maintainability"><img src="https://api.codeclimate.com/v1/badges/957eaee7cf558abcf2d0/maintainability" /></a> ![Sponsor](https://img.shields.io/badge/Sponsor-❤-ff69b4.svg)
</div>

### 💫 Features

> [!WARNING]
>
>API is currently under development. Expect potential changes and unstable behavior, but it is fully functional and ready for use.

- **Zero Dependencies**: Lightweight with no external dependencies.
- **Multi-Service Support**: Send notifications to Slack, Discord, Resend, and more.
- **Bulk Messaging**: Send notifications to multiple repositories simultaneously.
- **Extensible**: Easily add more services or custom logic.

### 📦 Installation

#### Install

```sh
go get -u github.com/lucasvillarinho/nofy
```

#### Example

##### Slack

```go
// Create a new Slack messenger
slackMessenger, _ := slack.NewSlackMensseger(
    // Set the Slack token to be used to send (required)
    slack.WithToken("test-token"),
    slack.WithMessage(
        // Message to be sent to the slack channel (required)
        // The message is a slice of maps, each map represents a block of the message
        // In this case, we are sending a single block with a text section
        slack.Message{
            Channel: "test-channel",
            Content: []map[string]any{
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
_ = nofy.SendAll(context.Background())
```

### 💛 Support the author

[![Sponsor](https://img.shields.io/badge/Sponsor-❤-ff69b4.svg)](https://github.com/sponsors/lucasvillarinho)
Enjoying the project? Consider supporting it to help me keep improving and adding new features!

[🙌 Support here](https://github.com/sponsors/lucasvillarinho)

### 🤝 Alternatives

For more example please check the specification file.

- [nikoksr/notify](https://github.com/nikoksr/notify)
- [containrrr/shoutrrr](https://github.com/containrrr/shoutrrr)
- [caronc/apprise](https://github.com/caronc/apprise)

### 📜 License

 ![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)

This software is licensed under the [MIT](https://github.com/lucasvillarinho/nofy/blob/main/LICENSE)
