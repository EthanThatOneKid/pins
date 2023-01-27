package discord

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/pins/preprocess"
)

// ChannelID is a Discord channel ID.
type ChannelID = discord.ChannelID

// client is a Discord client.
type client struct {
	*api.Client
}

// New creates a new Discord client.
func New(token string) preprocess.ChatProvider {
	return &client{api.NewClient(token)}
}

// Pins reads the pinned messages from the guild.
func (c *client) Pins(options preprocess.PinsOptions) (preprocess.PinsList, error) {
	guildID, err := discord.ParseSnowflake(options.GuildID)
	if err != nil {
		return preprocess.PinsList{}, err
	}

	guild, err := c.Client.Guild(discord.GuildID(guildID))
	if err != nil {
		return preprocess.PinsList{}, err
	}

	channels, err := c.Channels(options.GuildID)
	if err != nil {
		return preprocess.PinsList{}, err
	}

	// Get the pinned messages from the channels.
	pinsList := preprocess.PinsList{
		GuildID:   guild.ID.String(),
		GuildName: guild.Name,
		Pins:      []preprocess.Pins{},
	}
	var pins preprocess.Pins
	for _, channel := range channels {
		// Get the pinned messages from the channel.
		channelID, err := discord.ParseSnowflake(channel.ID)
		if err != nil {
			return preprocess.PinsList{}, err
		}

		messages, err := c.Client.PinnedMessages(discord.ChannelID(channelID))
		if err != nil {
			return preprocess.PinsList{}, err
		}

		pins.Messages = make([]preprocess.Message, len(messages))
		for i, message := range messages {
			pins.Messages[i] = convertMessage(message)
		}

		pinsList.Pins = append(pinsList.Pins, pins)
	}

	return pinsList, nil
}

// Channels gets all the channel IDs from a guild.
func (c *client) Channels(id string) ([]preprocess.Channel, error) {
	guildID, err := discord.ParseSnowflake(id)
	if err != nil {
		return nil, err
	}

	// Get all the text channels from the guild.
	everyChannel, err := c.Client.Channels(discord.GuildID(guildID))
	if err != nil {
		return nil, err
	}

	// Filter out the non-text channels.
	channels := []discord.Channel{}
	for _, channel := range everyChannel {
		if channel.Type == discord.GuildText {
			channels = append(channels, channel)
		}
	}
	channelsLen := len(channels)

	// Append archived threads to the channels.
	for i := 0; i < channelsLen; i++ {
		publicArchivedThreads, err := c.PublicArchivedThreads(channels[i].ID, discord.Timestamp{}, 0)
		if err != nil {
			return nil, err
		}
		channels = append(channels, publicArchivedThreads.Threads...)

		privateArchivedThreads, err := c.PrivateArchivedThreads(channels[i].ID, discord.Timestamp{}, 0)
		if err != nil {
			return nil, err
		}
		channels = append(channels, privateArchivedThreads.Threads...)
	}

	// Append active threads to the channels.
	activeThreads, err := c.ActiveThreads(discord.GuildID(guildID))
	if err != nil {
		return nil, err
	}
	channels = append(channels, activeThreads.Threads...)

	// Convert the channels to the preprocess format.
	return convertChannels(channels), nil
}

func convertMessage(message discord.Message) preprocess.Message {
	var attachments []preprocess.Attachment
	for _, attachment := range message.Attachments {
		attachments = append(attachments, convertAttachment(attachment))
	}

	return preprocess.Message{
		Timestamp:   message.Timestamp.Time(),
		AuthorID:    message.Author.ID.String(),
		Text:        message.Content,
		Attachments: attachments,
	}
}

func convertAttachment(attachment discord.Attachment) preprocess.Attachment {
	return preprocess.Attachment{
		URL:         attachment.URL,
		Filename:    attachment.Filename,
		Size:        attachment.Size,
		ContentType: attachment.ContentType,
	}
}

func convertChannels(channels []discord.Channel) []preprocess.Channel {
	var converted []preprocess.Channel
	for _, channel := range channels {
		converted = append(converted, preprocess.Channel{
			ID:   channel.ID.String(),
			Name: channel.Name,
		})
	}

	return converted
}
