package main

import (
	"fmt"
	"strings"

	"github.com/randallmlough/code-bot/internal/file"
	"github.com/randallmlough/code-bot/internal/openai"
	"github.com/urfave/cli/v2"
)

func startChat(c *cli.Context) error {
	// init a chat session
	sesh := openai.ContextGetSession(c.Context)
	for _, msg := range openai.PluginsChatConversations(sesh.Plugins()...) {
		sesh.AddMessage(msg)
	}

	inputs := c.StringSlice("i")
	if len(inputs) > 0 {
		for _, iSrc := range c.StringSlice("i") {
			input, err := file.ReadFile(iSrc)
			if err != nil {
				return fmt.Errorf("failed to read the input flag: %w", err)
			}
			sesh.AddUserMessage(input)
			printUserMessage(input)
		}
	}

	var prompt string
	args := c.Args()
	if args.Present() {
		prompt = strings.Join(args.Slice(), " ")
	}

	if c.String("p") != "" {
		var err error
		prompt, err = file.ReadFile(c.String("p"))
		if err != nil {
			return fmt.Errorf("failed to read the prompt flag: %w", err)
		}
	}

	if prompt == "" && len(inputs) > 0 {
		printBotMessage("What would you like me to do with these inputs?")
		prompt = waitForInput()
		if prompt == "" {
			printBotMessage("I need to know what to do with that input.")
			printErrorMessage("Exiting.")
			return nil
		}
	}

	// send request
	quit := false
	for !quit {
		if prompt == "" {
			printBotMessage("What can I help you with? (type `quit` or `exit` to close)")
			prompt = waitForInput()
		}
		// questionParam := validateQuestion(question)
		switch prompt {
		case "quit", "exit":
			quit = true
		case "":
			continue
		default:
			sesh.AddUserMessage(prompt)
			printUserMessage(prompt)

			printSystemMessage("handling request. This can take some time depending on the request.")
			ctx := c.Context
			resp, err := sesh.Send(ctx)
			if err != nil {
				return err
			}
			printBotMessage(resp.Message.Content)
		}
		prompt = ""
	}

	return nil
}

// -p "testdata/prompt.txt" -i "testdata/input.ts"
