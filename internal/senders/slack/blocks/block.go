package blocks

// BlockMessage is a struct that represents a block message in a Slack message.
type BlockMessage struct {
	Channel string  `json:"channel"`
	Blocks  []Block `json:"blocks"`
}

// Block is a struct that represents a block in a Slack message.
// It can be a section, context, actions, or a divider.
// Look at https://api.slack.com/reference/messaging/blocks for more information.
type Block struct {
	Type  string `json:"type"`
	Image Image  `json:"image,omitempty"`
}
