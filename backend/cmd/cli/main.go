package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/randallmlough/code-bot/internal/env"
	"github.com/randallmlough/code-bot/internal/openai"
	"github.com/randallmlough/code-bot/internal/openai/bots"
	projectassistant "github.com/randallmlough/code-bot/internal/openai/plugins/project_assistant"
	"github.com/randallmlough/code-bot/internal/tracedlogging"
	"github.com/segmentio/ksuid"
	"github.com/urfave/cli/v2"
)

func init() {
	h := slog.NewTextHandler(os.Stderr, nil)
	ctxH := tracedlogging.NewContextHandler(h)
	logger := slog.New(ctxH)
	slog.SetDefault(logger)
}

func main() {

	t := env.GetString("OPENAI_TOKEN", os.Getenv("OPENAI_TOKEN"))
	if t == "" {
		printErrorMessage("Missing OPENAI_TOKEN")
		os.Exit(1)
	}
	client := openai.New(t)
	botCfg := openai.BotConfig{
		FuncCallLimit: 4,
	}
	pa := projectassistant.New()
	plugins := []openai.Plugin{
		pa,
	}
	sesh, err := bots.NewCodingBot(client, botCfg, plugins...)
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	ctx = tracedlogging.SetTraceID(ctx, ksuid.New().String())
	if err := runCli(openai.ContextSetSession(ctx, sesh)); err != nil {
		printErrorMessage(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func runCli(ctx context.Context) error {
	app := &cli.App{
		Name:  "code-bot",
		Usage: "a helpful coding assistant",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "input",
			},
			&cli.StringFlag{
				Name:    "prompt",
				Aliases: []string{"p"},
				Usage:   "the prompt for what to do with the input",
			},
		},
		DefaultCommand: "chat",
		Commands: []*cli.Command{
			{
				Name:        "chat",
				Aliases:     []string{"c"},
				UsageText:   "code-bot chat [command options] [prompt...]",
				Description: "Starts a chat conversation",
				ArgsUsage:   "",
				Action:      startChat,
				Flags:       nil,
			},
			{
				Name:        "plugins",
				Aliases:     nil,
				Usage:       "",
				UsageText:   "",
				Description: "Manage plugins",
				ArgsUsage:   "",
				Action: func(c *cli.Context) error {
					fmt.Println(c.Args())
					sesh := openai.ContextGetSession(c.Context)

					fmt.Println("PLUGINS ACTION:", sesh.Plugins())
					return nil
				},
				Flags: nil,
			},
		},
	}
	return app.RunContext(ctx, os.Args)
}
