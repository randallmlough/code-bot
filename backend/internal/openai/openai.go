package openai

import (
	gogpt "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	c *gogpt.Client
}

func New(token string) *OpenAI {
	c := gogpt.NewClient(token)
	return &OpenAI{c}
}

type Conversation interface {
	Messages() []gogpt.ChatCompletionMessage
}
