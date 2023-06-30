package openai

import (
	"maps"

	gogpt "github.com/sashabaranov/go-openai"
)

// JsonObject is just a more readable type
type JsonObject = map[string]any

type Plugin interface {
	Name() string
	Description() string
	Parameters() JSONSchemaDefinition
	Execute(args JsonObject) (JsonObject, error)
}

// PluginPrompt will inject additional context into the system prompt
type PluginPrompt interface {
	SystemPrompt() string
}

type Plugins map[string]Plugin

func (pp Plugins) FunctionDefinitions() []FunctionDefinition {
	fns := []FunctionDefinition{}
	plugins := maps.Values(pp)
	for _, plugin := range plugins {
		fn := FunctionDefinition{
			Name:        plugin.Name(),
			Description: plugin.Description(),
			Parameters:  plugin.Parameters(),
		}
		fns = append(fns, fn)
	}
	return fns
}

func (pp Plugins) GetPlugin(name string) Plugin {
	plugin, exists := pp[name]
	if !exists {
		return nil
	}
	return plugin
}

// PluginsChatConversations is a helper function that will
func PluginsChatConversations(plugins ...Plugin) []gogpt.ChatCompletionMessage {
	messages := []gogpt.ChatCompletionMessage{}
	for _, plugin := range plugins {
		if p, ok := plugin.(Conversation); ok {
			messages = append(messages, p.Messages()...)
		}
	}
	return messages
}
