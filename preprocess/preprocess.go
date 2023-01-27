package preprocess

import (
	"regexp"
	"time"
)

// ChatProvider provides a channel by ID from any service.
type ChatProvider interface {
	// Pins reads the pins of a guild by guild ID.
	Pins(o PinsOptions) (PinsList, error)
	// Channels gets all the channels from a guild including threads.
	Channels(id string) ([]Channel, error)
}

type PinsOptions struct {
	// GuildID is the ID of the guild to get the pins from.
	GuildID string
	// KeepMessages is whether to keep the messages in the pins.
	KeepMessages []regexp.Regexp
	// SkipMessages is whether to skip the messages in the pins.
	SkipMessages []regexp.Regexp
	// KeepChannels is whether to keep the channels in the pins.
	KeepChannels []regexp.Regexp
	// SkipChannels is whether to skip the channels in the pins.
	SkipChannels []regexp.Regexp
}

type Channel struct {
	ID   string
	Name string
}

type MessageList struct {
	Messages []Message
	Authors  map[string]Author
}

type Pins struct {
	Channel
	MessageList
}

type PinsList struct {
	GuildID   string
	GuildName string
	Pins      []Pins
}

// Message is a text message struct.
type Message struct {
	// Timestamp is the time the message was sent.
	Timestamp time.Time `json:"timestamp"`
	// AuthorID is the user who sent the message.
	AuthorID string `json:"author_id"`
	// Text is the content of the message.
	Text string `json:"text"`
	// Attachments is a list of attachments.
	Attachments []Attachment `json:"attachments"`
}

// Author is a user struct.
type Author struct {
	// Name is the name of the user.
	Name string `json:"name"`
	// ID is the unique identifier of the user.
	ID string `json:"id"`
	// Avatar is the URL of the user's avatar.
	Avatar string `json:"avatar"`
}

// Attachment is a file attachment.
type Attachment struct {
	// URL is the URL of the attachment.
	URL string `json:"url"`
	// Filename is the name of the file.
	Filename string `json:"filename"`
	// Size is the size of the file in bytes.
	Size uint64 `json:"size"`
	// ContentType is the media type of the file.
	ContentType string `json:"content_type"`
}
