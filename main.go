package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	arikawa "github.com/diamondburned/arikawa/v3/discord"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ethanthatonekid/pins/pins"
	"github.com/ethanthatonekid/pins/pins/discord"
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
		Name:      "pins",
		Usage:     "get pins from Discord",
		ArgsUsage: "guildID [expression]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "print information about each pin",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "don't save pins",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("verbose") {
				log.SetOutput(ctx.App.ErrWriter)
				log.SetFlags(log.Ltime)
			} else {
				log.SetOutput(io.Discard)
			}
			return nil
		},
		After: func(ctx *cli.Context) error {
			log.SetOutput(ctx.App.ErrWriter)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "get",
				Action:    getAction,
				Usage:     "get pins from a guild",
				ArgsUsage: "guildID [expression]",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "output",
						Usage:   "the output directory",
					},
				},
			},
			{
				Name:      "update",
				Action:    updateAction,
				Usage:     "update the output file with new pins",
				ArgsUsage: "file",
			},
			{
				Name:      "print-expr",
				Action:    printExprAction,
				Usage:     "print the expression used to run this tool from a path",
				ArgsUsage: "filepath",
			},
		},
	}

	return app
}

func getAction(ctx *cli.Context) error {
	provider := discord.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	arg1 := ctx.Args().Get(0)
	expr := pins.FilterExpr(ctx.Args().Get(1))

	snowflake, err := arikawa.ParseSnowflake(arg1)
	if err != nil {
		return errors.Wrap(err, "failed to parse snowflake")
	}

	guildID := arikawa.GuildID(snowflake)

	dryRun := ctx.Bool("dry-run")
	outputRoot := ctx.Path("output")

	data, err := provider.Pins(pins.PinsOptions{
		GuildID: guildID,
		Expr:    expr,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get pins")
	}

	log.Println(data)

	if dryRun {
		return nil
	}

	if err := os.MkdirAll(outputRoot, 0755); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	filename, err := encodeFilename(guildID, expr)
	if err != nil {
		return errors.Wrap(err, "failed to encode filename")
	}

	outputPath := filepath.Join(outputRoot, filename)

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return errors.Wrap(err, "failed to write output file")
	}

	return nil
}

func printExprAction(ctx *cli.Context) error {
	path := ctx.Args().Get(0)
	if path == "" {
		return errors.New("no file provided, check -h")
	}

	guildID, hash, err := decodeFilename(filepath.Base(path))
	if err != nil {
		return errors.Wrap(err, "failed to decode filename")
	}

	fmt.Println("Copy the following flags to run this tool:")
	fmt.Printf("  --guild %q --expr %q\n", guildID, hash)

	return nil
}

func updateAction(ctx *cli.Context) error {
	provider := discord.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	path := ctx.Args().Get(0)
	if path == "" {
		return errors.New("no file provided, check -h")
	}

	guildID, expr, err := decodeFilename(filepath.Base(path))
	if err != nil {
		return errors.Wrap(err, "failed to decode filename")
	}

	data, err := provider.Pins(pins.PinsOptions{
		GuildID: guildID,
		Expr:    expr,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get pins")
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return errors.Wrap(err, "failed to write output file")
	}

	return nil
}
