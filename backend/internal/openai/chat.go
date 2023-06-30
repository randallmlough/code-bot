package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"

	gogpt "github.com/sashabaranov/go-openai"
	"github.com/segmentio/ksuid"
)

type Chat struct {
	c       *gogpt.Client
	id      string
	plugins Plugins

	messages       []gogpt.ChatCompletionMessage
	sendIterations int
	cfg            *BotConfig
}

type BotConfig struct {
	Model            string         `json:"model,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	FuncCallLimit    int            `json:"-"`
}

func NewChatSession(client *OpenAI, opts ...Option) *Chat {
	srv := &Chat{
		c:        client.c,
		id:       ksuid.New().String(),
		plugins:  map[string]Plugin{},
		messages: []gogpt.ChatCompletionMessage{},
	}
	if len(opts) > 0 {
		srv = srv.WithOptions(opts...)
	}
	return srv
}

func (c *Chat) clone() *Chat {
	srvCopy := *c
	return &srvCopy
}

// An Option configures a Logger.
type Option interface {
	apply(*Chat)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Chat)

func (f optionFunc) apply(srv *Chat) {
	f(srv)
}

// WithOptions clones the current Service, applies the supplied Options, and
// returns the resulting Service.
func (c *Chat) WithOptions(opts ...Option) *Chat {
	s := c.clone()
	for _, opt := range opts {
		opt.apply(s)
	}
	return s
}

func WithPlugins(plugins []Plugin) Option {
	return optionFunc(func(s *Chat) {
		for _, plugin := range plugins {
			s.plugins[plugin.Name()] = plugin
		}
	})
}

func SetConfig(cfg BotConfig) Option {
	return optionFunc(func(s *Chat) {
		s.cfg = &cfg
	})
}

func WithMessages(msgs []gogpt.ChatCompletionMessage) Option {
	return optionFunc(func(s *Chat) {
		s.messages = append(s.messages, msgs...)
	})
}

func (c *Chat) RegisterPlugin(plugin Plugin) {
	c.plugins[plugin.Name()] = plugin
}

func (c *Chat) Plugins() []Plugin {
	return maps.Values(c.plugins)
}

func (c *Chat) GetPlugin(name string) Plugin {
	return c.plugins[name]
}

func NewMessage(role, message string) gogpt.ChatCompletionMessage {
	return gogpt.ChatCompletionMessage{
		Role:    role,
		Content: message,
	}
}

func NewUserMessage(message string) gogpt.ChatCompletionMessage {
	return gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: message,
	}
}

func NewAssistantMessage(message string) gogpt.ChatCompletionMessage {
	return gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleAssistant,
		Content: message,
	}
}

func NewFunctionMessage(message, functionName string) gogpt.ChatCompletionMessage {
	return gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleFunction,
		Content: message,
		Name:    functionName,
	}
}

func (c *Chat) AddSystemMessage(message string) {
	c.addMessage(gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleSystem,
		Content: message,
	})
}

func (c *Chat) AddUserMessage(message string) {
	c.addMessage(NewUserMessage(message))
}

func (c *Chat) AddAssistantMessage(message gogpt.ChatCompletionMessage) {
	c.addMessage(message)
}

func (c *Chat) AddFunctionMessage(message, functionName string) {
	c.addMessage(gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleFunction,
		Content: message,
		Name:    functionName,
	})
}

func (c *Chat) AddMessage(msg gogpt.ChatCompletionMessage) {
	c.addMessage(msg)
}

// addMessage executes the message
func (c *Chat) addMessage(msg gogpt.ChatCompletionMessage) {
	c.messages = append(c.messages, msg)
}

// Send sends the request to OpenAI and returns back a response
func (c *Chat) Send(ctx context.Context) (resp gogpt.ChatCompletionChoice, err error) {
	c.clearRequestIteration()
	return c.send(ctx)
}

// Send sends the request to OpenAI and returns back a response
func (c *Chat) send(ctx context.Context) (gogpt.ChatCompletionChoice, error) {
	req, bErr := c.buildRequest(c.cfg)
	if bErr != nil {
		return gogpt.ChatCompletionChoice{}, bErr
	}

	r, err := c.c.CreateChatCompletion(ctx, req)
	if err != nil {
		return gogpt.ChatCompletionChoice{}, err
	}
	resp := r.Choices[0]

	c.AddAssistantMessage(resp.Message)

	if resp.Message.FunctionCall != nil {
		if c.reachedFuncCallLimit() {
			return resp, errors.New("reached function call limit")
		}
		plugin := c.plugins.GetPlugin(resp.Message.FunctionCall.Name)
		args := map[string]any{}
		if jsonErr := json.Unmarshal([]byte(resp.Message.FunctionCall.Arguments), &args); jsonErr != nil {
			return resp, fmt.Errorf("failed to unmarshal args to function call %q: %w", plugin.Name(), jsonErr)
		}
		pluginResponse, pErr := plugin.Execute(args)
		if pErr != nil {
			return resp, fmt.Errorf("failed exeucting plugin: %q err: %w", plugin.Name(), pErr)
		}
		b, jsonErr := json.Marshal(pluginResponse)
		if jsonErr != nil {
			return resp, fmt.Errorf("failed marshall response for plugin: %q err: %w", plugin.Name(), jsonErr)
		}
		c.AddFunctionMessage(string(b), plugin.Name())
		c.increaseRequestIteration()
		return c.send(ctx)
	}
	return resp, err
}

func (c *Chat) clearRequestIteration() {
	c.sendIterations = 0
}

func (c *Chat) increaseRequestIteration() {
	c.sendIterations += 1
}

func (c *Chat) reachedFuncCallLimit() bool {
	// zero and below is treated as unlimited
	if c.cfg.FuncCallLimit <= 0 {
		return false
	}
	return c.sendIterations > c.cfg.FuncCallLimit
}

type SendOptions struct {
	FuncCallLimit int
}

// SendWithOptions sends the request to OpenAI
func (c *Chat) SendWithOptions(ctx context.Context, opts *SendOptions) (gogpt.ChatCompletionChoice, error) {
	if opts != nil {

	}
	return c.Send(ctx)
}

var defaultRequestConfig = gogpt.ChatCompletionRequest{
	Model:       gogpt.GPT3Dot5Turbo0613,
	Temperature: 0.7,
}

func (c *Chat) buildRequest(cfg *BotConfig) (gogpt.ChatCompletionRequest, error) {
	req := defaultRequestConfig
	if len(c.messages) > 0 {
		req.Messages = c.messages
	}

	if len(c.plugins) > 0 {
		req.Functions = c.plugins.FunctionDefinitions()
	}
	if cfg != nil {
		b, err := json.Marshal(cfg)
		if err != nil {
			return req, err
		}
		if err := json.Unmarshal(b, &req); err != nil {
			return req, err
		}
	}
	return req, nil
}

func (c *Chat) GetMessage() {
}

func (c *Chat) GetMessages() {
}
