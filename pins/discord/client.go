package discord

import (
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/pins/pins"
	"github.com/ethanthatonekid/pins/pins/filter"
)

// client is a Discord client.
type client struct {
	*api.Client
}

// New creates a new Discord client.
func New(token string) pins.ChatProvider {
	return &client{api.NewClient(token)}
}

// Pins reads the pinned messages from the guild.
func (c *client) Pins(options pins.PinsOptions) (*pins.Pins, error) {
	ftr, err := filter.New(options.Expr)
	if err != nil {
		return nil, err
	}

	guild, err := c.Client.Guild(options.GuildID)
	if err != nil {
		return nil, err
	}

	log.Println("working on guild", guild.Name)

	channels, err := c.Channels(options.GuildID)
	if err != nil {
		return nil, err
	}

	result := &pins.Pins{
		GuildID:      guild.ID,
		GuildName:    guild.Name,
		ChannelNames: map[discord.ChannelID]string{},
		Authors:      map[discord.UserID]pins.Author{},
		Pins:         []pins.Pin{},
	}

	log.Println("found", len(channels), "channels")

	for _, channel := range channels {
		data := filter.MessageData{
			ChannelParentID: channel.ParentID,
			ChannelID:       channel.ID,
			ChannelName:     channel.Name,
		}
		keep, err := ftr.Filter(data)
		if err != nil {
			return nil, err
		}

		if keep {
			log.Println("keeping channel", channel.Name)
		} else {
			log.Println("ignoring channel", channel.Name)
			continue
		}

		result.ChannelNames[channel.ID] = channel.Name

		messages, err := c.Client.PinnedMessages(channel.ID)
		if err != nil {
			return nil, err
		}

		for _, message := range messages {
			data.AuthorID = message.Author.ID
			data.AuthorName = message.Author.Username
			data.Timestamp = message.Timestamp.Time()
			data.Text = message.Content
			keep, err := ftr.Filter(data)
			if err != nil {
				return nil, err
			}

			// log.Println("message", message.Author.Tag(), keep)

			if !keep {
				continue
			}

			result.Pins = append(result.Pins, pins.Pin{
				ChannelID:   channel.ID,
				Timestamp:   message.Timestamp.Time(),
				AuthorID:    message.Author.ID,
				Text:        message.Content,
				Attachments: convertAttachments(message.Attachments),
			})

			mentions := []pins.Author{{
				ID:     message.Author.ID,
				Name:   message.Author.Tag(),
				Avatar: message.Author.AvatarURL(),
			}}
			for _, mentioned := range message.Mentions {
				mentions = append(mentions, pins.Author{
					ID:     mentioned.ID,
					Name:   mentioned.Tag(),
					Avatar: mentioned.AvatarURL(),
				})
			}

			for _, mentioned := range mentions {
				if _, ok := result.Authors[mentioned.ID]; !ok {
					result.Authors[mentioned.ID] = mentioned
				}
			}
		}
	}

	return result, nil
}

// Channels gets all the channel IDs from a guild.
func (c *client) Channels(guildID discord.GuildID) ([]pins.Channel, error) {
	// Get all the text channels from the guild.
	everyChannel, err := c.Client.Channels(guildID)
	if err != nil {
		return nil, err
	}

	// Filter out the non-text channels.
	channels := everyChannel[:0]
	for _, channel := range everyChannel {
		if channel.Type == discord.GuildText {
			channels = append(channels, channel)
		}
	}

	// Append archived threads to the channels.
	for i := range channels {
		log.Println("fetching public archived threads for channel", channels[i].Name)
		publicArchivedThreads, err := c.PublicArchivedThreads(channels[i].ID, discord.Timestamp{}, 0)
		if err != nil {
			return nil, err
		}
		channels = append(channels, publicArchivedThreads.Threads...)

		log.Println("fetching private archived threads for channel", channels[i].Name)
		privateArchivedThreads, err := c.PrivateArchivedThreads(channels[i].ID, discord.Timestamp{}, 0)
		if err != nil {
			return nil, err
		}
		channels = append(channels, privateArchivedThreads.Threads...)
	}

	log.Println("fetching active threads")
	// Append active threads to the channels.
	activeThreads, err := c.ActiveThreads(discord.GuildID(guildID))
	if err != nil {
		return nil, err
	}
	channels = append(channels, activeThreads.Threads...)

	// Convert the channels to the pins format.
	return convertChannels(channels), nil
}

func convertAttachments(attachment []discord.Attachment) []pins.Attachment {
	var converted []pins.Attachment
	for _, attachment := range attachment {
		converted = append(converted, convertAttachment(attachment))
	}
	return converted
}

func convertAttachment(attachment discord.Attachment) pins.Attachment {
	return pins.Attachment{
		URL:         attachment.URL,
		Filename:    attachment.Filename,
		Size:        attachment.Size,
		ContentType: attachment.ContentType,
	}
}

func convertChannels(channels []discord.Channel) []pins.Channel {
	var converted []pins.Channel
	for _, channel := range channels {
		converted = append(converted, pins.Channel{
			ID:       channel.ID,
			Name:     channel.Name,
			ParentID: channel.ParentID,
		})
	}

	return converted
}
