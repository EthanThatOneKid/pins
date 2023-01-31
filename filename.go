package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/pins/pins"
)

func encodeFilename(id discord.GuildID, expr pins.FilterExpr) (string, error) {
	if expr == "" {
		return fmt.Sprintf("%s.json", id), nil
	}

	return fmt.Sprintf(
		"%s-%s.json",
		id,
		base64.URLEncoding.EncodeToString([]byte(expr)),
	), nil
}

func decodeFilename(filename string) (discord.GuildID, pins.FilterExpr, error) {
	filename = strings.TrimSuffix(filename, ".json")

	var expr pins.FilterExpr

	snowflake, encoded, found := strings.Cut(filename, "-")
	if found {
		exprBytes, err := base64.URLEncoding.DecodeString(encoded)
		if err != nil {
			return 0, "", fmt.Errorf("base64 error: %w", err)
		}
		expr = pins.FilterExpr(exprBytes)
	} else {
		snowflake = filename
	}

	guildID, err := discord.ParseSnowflake(snowflake)
	if err != nil {
		return 0, "", fmt.Errorf("cannot parse snowflake: %w", err)
	}

	return discord.GuildID(guildID), expr, nil
}
