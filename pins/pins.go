package pins

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

// ChatProvider provides a channel by ID from any service.
type ChatProvider interface {
	// Pins reads the pins of a guild by guild ID.
	Pins(PinsOptions) (*Pins, error)
	// Channels gets all the channels from a guild including threads.
	Channels(guildID discord.GuildID) ([]Channel, error)
}

// FilterExpr is a CEL expression for filtering.
type FilterExpr string

// PinsOptions are the options for getting pins.
type PinsOptions struct {
	// GuildID is the ID of the guild to get the pins from.
	GuildID discord.GuildID
	// CEL expression to filter pins.
	Expr FilterExpr
}

type Pins struct {
	// GuildID is the ID of the guild the pins are from.
	GuildID discord.GuildID `json:"guild_id"`
	// GuildName is the name of the guild the pins are from.
	GuildName string `json:"guild_name"`
	// ChannelNames is a map of channel IDs to names.
	ChannelNames map[discord.ChannelID]string `json:"channel_names"`
	// Authors is a map of user IDs to names.
	Authors map[discord.UserID]Author `json:"authors"`
	// Pins is a list of pins.
	Pins []Pin `json:"pins"`
}

// Pin is a pinned message struct.
type Pin struct {
	// ChannelID is the channel the message was sent in.
	ChannelID discord.ChannelID `json:"channel_id"`
	// Timestamp is the time the message was sent.
	Timestamp time.Time `json:"timestamp"`
	// AuthorID is the user who sent the message.
	AuthorID discord.UserID `json:"author_id"`
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
	ID discord.UserID `json:"id"`
	// Avatar is the URL of the user's avatar.
	Avatar string `json:"avatar"`
}

// Attachment is a file attachment.
type Attachment = discord.Attachment

// Channel is a channel struct.
type Channel struct {
	// ID is the unique identifier of the channel.
	ID discord.ChannelID `json:"id"`
	// Name is the name of the channel.
	Name string `json:"name"`
	// ParentID is the ID of the parent channel or category.
	ParentID discord.ChannelID `json:"parent_id"`
}
