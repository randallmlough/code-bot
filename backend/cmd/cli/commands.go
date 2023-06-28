package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/randallmlough/code-bot/internal/env"
	"github.com/randallmlough/code-bot/internal/file"
	"github.com/randallmlough/code-bot/internal/openai"
	"github.com/randallmlough/code-bot/internal/openai/bots"
	projectassistant "github.com/randallmlough/code-bot/internal/openai/plugins/project_assistant"
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	ctx := c.Context
	// init a chat session
	t := env.GetString("OPENAI_TOKEN", os.Getenv("OPENAI_TOKEN"))
	client := openai.New(t)
	botCfg := openai.BotConfig{}
	pa := projectassistant.New()
	plugins := []openai.Plugin{
		pa,
	}
	sesh, err := bots.NewCodingBot(client, botCfg, plugins...)
	if err != nil {
		return err
	}
	sesh.AddUserMessage(pa.UserMessage())
	sesh.AddMessage(gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleAssistant,
		Content: "I have stored your project structure. What can I help you with?",
	})
	args := c.Args()
	var haveInitialRequest bool
	if args.Present() {
		haveInitialRequest = true
		prompt := strings.Join(args.Slice(), " ")
		if err != nil {
			return fmt.Errorf("failed to read the prompt flag: %w", err)
		}
		sesh.AddUserMessage(prompt)
		printUserMessage(prompt)
	}

	// Send request
	if c.String("p") != "" {
		haveInitialRequest = true
		prompt, err := file.ReadFile(c.String("p"))
		if err != nil {
			return fmt.Errorf("failed to read the prompt flag: %w", err)
		}
		sesh.AddUserMessage(prompt)
		printUserMessage(prompt)
	}

	if c.String("i") != "" {
		if !haveInitialRequest {

		}
		input, err := file.ReadFile(c.String("i"))
		if err != nil {
			return fmt.Errorf("failed to read the input flag: %w", err)
		}
		sesh.AddUserMessage(input)
		printUserMessage(input)
	}

	// check if there was any input or flags

	// loop

	// send request
	printSystemMessage("handling request. This can take some time depending on the request.")
	resp, err := sesh.SendAndHandle(ctx)
	if err != nil {
		panic(err)
	}
	printBotMessage(resp.Message.Content)
	return nil
}
