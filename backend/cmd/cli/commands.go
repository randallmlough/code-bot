package main

import (
	"fmt"
	"os"

	"github.com/randallmlough/code-bot/internal/env"
	"github.com/randallmlough/code-bot/internal/openai"
	"github.com/randallmlough/code-bot/internal/openai/bots"
	projectassistant "github.com/randallmlough/code-bot/internal/openai/plugins/project_assistant"
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func Action(c *cli.Context) error {
	ctx := c.Context
	// get inputs
	prompt := c.String("p")
	input := c.String("i")

	// create prompt/message
	constructedPrompt := fmt.Sprintf("%s\ninput: ###\n%s\n###", prompt, input)
	// constructedPrompt := prompt
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
	// add message
	slog.Info("convo", "user", constructedPrompt)
	sesh.AddUserMessage(constructedPrompt)

	// send request
	slog.Info("convo", "system", "handling request. This can take some time depending on the request.")
	resp, err := sesh.SendAndHandle(ctx)
	if err != nil {
		panic(err)
	}
	slog.Info("convo", "bot", resp.Message.Content)
	return nil
}
