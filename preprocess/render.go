package preprocess

import (
	"regexp"

	"github.com/pkg/errors"
)

// RenderedGuild is a guild's pins rendered as JSON.
type RenderedGuild struct {
	GuildID   string             `json:"guild_id"`
	GuildName string             `json:"guild_name"`
	Channels  []*RenderedChannel `json:"channels"`
}

// RenderedChannel is a channel's messages rendered as JSON.
type RenderedChannel struct {
	ChannelID   string                    `json:"channel_id"`
	ChannelName string                    `json:"channel_name"`
	Authors     map[string]RenderedAuthor `json:"authors"`
	Messages    []RenderedMessage         `json:"messages"`
}

type RenderedMessage = Message

type RenderedAuthor struct {
	// Name is the name of the user.
	Name string `json:"name"`
	// Avatar is the URL of the user's avatar.
	Avatar string `json:"avatar"`
}

type RenderGuildOptions struct {
	GuildID      string
	SkipChannels []string
	KeepChannels []string
	SkipMessages []string
	KeepMessages []string
}

// RenderPinsList renders the pins list of a guild to JSON.
func RenderGuild(provider ChatProvider, options RenderGuildOptions) (*RenderedGuild, error) {
	skipChannels, err := convertStringSliceToRegexpSlice(options.SkipChannels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert skip channels to regexps")
	}

	keepChannels, err := convertStringSliceToRegexpSlice(options.KeepChannels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert keep channels to regexps")
	}

	skipMessages, err := convertStringSliceToRegexpSlice(options.SkipMessages)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert skip messages to regexps")
	}

	keepMessages, err := convertStringSliceToRegexpSlice(options.KeepMessages)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert keep messages to regexps")
	}

	pinsList, err := provider.Pins(PinsOptions{
		GuildID:      options.GuildID,
		SkipChannels: skipChannels,
		KeepChannels: keepChannels,
		SkipMessages: skipMessages,
		KeepMessages: keepMessages,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read pins")
	}

	var guild RenderedGuild
	guild.GuildID = pinsList.GuildID
	guild.GuildName = pinsList.GuildName

	for _, pins := range pinsList.Pins {
		renderedChannel, err := RenderChannel(pins)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render channel")
		}

		guild.Channels = append(guild.Channels, renderedChannel)
	}

	return &guild, nil
}

// RenderChannel renders a channel's pins to JSON.
func RenderChannel(pins Pins) (*RenderedChannel, error) {
	channel := &RenderedChannel{
		ChannelID:   pins.Channel.ID,
		ChannelName: pins.Channel.Name,
		Authors:     make(map[string]RenderedAuthor, len(pins.MessageList.Authors)),
		Messages:    make([]RenderedMessage, len(pins.MessageList.Messages)),
	}

	for i, msg := range pins.MessageList.Messages {
		channel.Messages[i] = RenderMessage(msg)
	}

	for id, author := range pins.MessageList.Authors {
		channel.Authors[id] = RenderAuthor(author)
	}

	return channel, nil
}

// RenderMessage renders a message to JSON.
func RenderMessage(msg Message) RenderedMessage {
	return msg
}

// RenderAuthor renders an author to JSON.
func RenderAuthor(author Author) RenderedAuthor {
	return RenderedAuthor{
		Name:   author.Name,
		Avatar: author.Avatar,
	}
}

func convertStringSliceToRegexpSlice(slice []string) ([]regexp.Regexp, error) {
	var regexps []regexp.Regexp
	for _, str := range slice {
		re, err := regexp.Compile(str)
		if err != nil {
			return nil, errors.Wrap(err, "failed to compile regexp")
		}

		regexps = append(regexps, *re)
	}

	return regexps, nil
}
