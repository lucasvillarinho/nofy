package blocks

type TextType string

const (
	PlainText TextType = "plain_text"
	Markdown  TextType = "mrkdwn"
)

// Text is a struct that represents a text in a Slack message.
type Text struct {
	Type  TextType `json:"type"`
	Text  string   `json:"text"`
	Emoji bool     `json:"emoji,omitempty"`
}
