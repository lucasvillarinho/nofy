package blocks

// Message represents a message to be sent to Slack.
type Message struct {
	Channel    string `json:"channel"`
	Text       string `json:"text"`
	IsMarkdown bool   `json:"mrkdwn"`
}
