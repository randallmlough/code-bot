package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	gogpt "github.com/sashabaranov/go-openai"
	"github.com/segmentio/ksuid"
)

type OpenAI struct {
	c *gogpt.Client
}

func New(token string) *OpenAI {
	c := gogpt.NewClient(token)
	return &OpenAI{c}
}

type Chat struct {
	c       *gogpt.Client
	id      string
	plugins Plugins

	messages []gogpt.ChatCompletionMessage
	cfg      *BotConfig
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
}

func (client *OpenAI) NewChatSession(opts ...Option) *Chat {
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

func (srv *Chat) clone() *Chat {
	srvCopy := *srv
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
func (srv *Chat) WithOptions(opts ...Option) *Chat {
	s := srv.clone()
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

func (s *Chat) RegisterPlugin(plugin Plugin) {
	s.plugins[plugin.Name()] = plugin
}

func (s *Chat) Plugins() []Plugin {
	return maps.Values(s.plugins)
}

func (s *Chat) AddSystemMessage(message string) {
	s.addMessage(gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleSystem,
		Content: message,
	})
}
func (s *Chat) AddUserMessage(message string) {
	s.addMessage(gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: message,
	})
}

func (s *Chat) AddAssistantMessage(message gogpt.ChatCompletionMessage) {
	s.addMessage(message)
}

func (s *Chat) AddFunctionMessage(message, functionName string) {
	s.addMessage(gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleFunction,
		Content: message,
		Name:    functionName,
	})
}

type ChatResponse struct {
	id string
}

func (s *Chat) AddMessage(msg gogpt.ChatCompletionMessage) {
	s.addMessage(msg)
}

// addMessage executes the message
func (s *Chat) addMessage(msg gogpt.ChatCompletionMessage) {
	s.messages = append(s.messages, msg)
}

// Send sends the request to OpenAI and returns back a response
func (s *Chat) Send(ctx context.Context) (gogpt.ChatCompletionChoice, error) {
	req, err := s.buildRequest(s.cfg)
	if err != nil {
		return gogpt.ChatCompletionChoice{}, err
	}

	r, err := s.c.CreateChatCompletion(ctx, req)
	if err != nil {
		return gogpt.ChatCompletionChoice{}, err
	}
	assistant := r.Choices[0]
	s.AddAssistantMessage(assistant.Message)
	return assistant, nil
}

// SendAndHandle sends the request to OpenAI and will call any necessary plugins if needed
func (s *Chat) SendAndHandle(ctx context.Context) (gogpt.ChatCompletionChoice, error) {
	resp, err := s.Send(ctx)
	if err != nil {
		return gogpt.ChatCompletionChoice{}, err
	}
	if resp.Message.FunctionCall != nil {
		plugin := s.plugins.GetPlugin(resp.Message.FunctionCall.Name)
		args := map[string]any{}
		if err := json.Unmarshal([]byte(resp.Message.FunctionCall.Arguments), &args); err != nil {
			return gogpt.ChatCompletionChoice{}, fmt.Errorf("failed to unmarshal args to function call %q: %w", plugin.Name(), err)
		}
		pluginResponse, err := plugin.Execute(args)
		if err != nil {
			return gogpt.ChatCompletionChoice{}, fmt.Errorf("failed exeucting plugin: %q err: %w", plugin.Name(), err)
		}
		b, err := json.Marshal(pluginResponse)
		if err != nil {
			return gogpt.ChatCompletionChoice{}, fmt.Errorf("failed marshall response for plugin: %q err: %w", plugin.Name(), err)
		}

		s.AddFunctionMessage(string(b), plugin.Name())
		return s.SendAndHandle(ctx)
	}
	// determine if a plugin call is needed to be made
	// loop through responses until no more plugins are needed
	return resp, nil
}

var defaultRequestConfig = gogpt.ChatCompletionRequest{
	Model:       gogpt.GPT3Dot5Turbo0613,
	Temperature: 0.7,
}

func (s *Chat) buildRequest(cfg *BotConfig) (gogpt.ChatCompletionRequest, error) {
	req := defaultRequestConfig
	if len(s.messages) > 0 {
		req.Messages = s.messages
	}

	if len(s.plugins) > 0 {
		req.Functions = s.plugins.FunctionDefinitions()
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

func (s *Chat) GetMessage() {
}

func (s *Chat) GetMessages() {
}
