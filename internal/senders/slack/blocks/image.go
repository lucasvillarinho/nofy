package blocks

type Image struct {
	Title    Text   `json:"title"`
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

type ImageOption func(image *Image)

func NewImageBlock(imageURL, altText string, options ...ImageOption) *Block {
	image := &Image{
		ImageURL: imageURL,
		AltText:  altText,
	}

	for _, opt := range options {
		opt(image)
	}

	return &Block{
		Type: "image",
		Image: Image{
			ImageURL: imageURL,
			AltText:  altText,
		},
	}
}

func WithTitle(text string) ImageOption {
	return func(s *Image) {
		s.Title = Text{
			Type:  "plain_text",
			Text:  text,
			Emoji: true,
		}
	}
}
