package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ethanthatonekid/pins/preprocess"
	"github.com/ethanthatonekid/pins/preprocess/discord"
)

func main() {
	godotenv.Load(".env")
	app := NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	*cli.App
}

func NewApp() *App {
	app := &App{}

	app.App = &cli.App{
		Name:     "pins",
		HelpName: "get pins from Discord",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "guild",
				Aliases: []string{"g"},
				Usage:   "the guild ID to save pins from",
			},
			&cli.StringSliceFlag{
				Name:    "skip-channel",
				Aliases: []string{"sc"},
				Usage:   "skipped if the channel contains the regex pattern",
			},
			&cli.StringSliceFlag{
				Name:    "keep-channel",
				Aliases: []string{"kc"},
				Usage:   "skipped if the channel does not contain the regex pattern",
			},
			&cli.StringSliceFlag{
				Name:    "skip-message",
				Aliases: []string{"sm"},
				Usage:   "skipped if the pin contains the regex pattern",
			},
			&cli.StringSliceFlag{
				Name:    "keep-message",
				Aliases: []string{"km"},
				Usage:   "skipped if the pin does not contain the regex pattern",
			},
			&cli.PathFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "output",
				Usage:   "the output directory",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "print information about each pin",
			},
		},
		Action: action,
	}

	return app
}

func action(ctx *cli.Context) error {
	provider := discord.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	outputRoot := ctx.Path("output")
	guildID := ctx.String("guild")
	verbose := ctx.Bool("verbose")

	data, err := preprocess.RenderGuild(provider, preprocess.RenderGuildOptions{
		GuildID:      guildID,
		SkipChannels: ctx.StringSlice("skip-channel"),
		KeepChannels: ctx.StringSlice("keep-channel"),
		SkipMessages: ctx.StringSlice("skip-message"),
		KeepMessages: ctx.StringSlice("keep-message"),
	})
	if err != nil {
		return errors.Wrap(err, "failed to render guild")
	}

	if verbose {
		fmt.Println(data)
	}

	if err := os.MkdirAll(outputRoot, 0755); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	outputFile := filepath.Join(outputRoot, guildID+".json")
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	if err := ioutil.WriteFile(outputFile, output, 0644); err != nil {
		return errors.Wrap(err, "failed to write output file")
	}

	return nil
}
