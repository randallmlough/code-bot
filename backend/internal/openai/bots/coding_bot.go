package bots

import (
	"embed"

	"github.com/randallmlough/code-bot/internal/file"
	"github.com/randallmlough/code-bot/internal/openai"
	codestyleassistant "github.com/randallmlough/code-bot/internal/openai/plugins/code_style_assistant"
)

//go:embed system_prompt.md
var files embed.FS

// create a plugin that gives examples of file structure, patterns, etc.

func NewCodingBot(client *openai.OpenAI, cfg openai.BotConfig, plugins ...openai.Plugin) (*openai.Chat, error) {

	sesh := openai.NewChatSession(
		client,
		openai.SetConfig(cfg),
		// add plugin that fetches latest version of a language
		openai.WithPlugins(plugins),
	)

	systemPrompt, err := file.ReadEmbededFile(files, "system_prompt.md")
	if err != nil {
		return nil, err
	}

	cs, err := codestyleassistant.New()
	if err != nil {
		return nil, err
	}
	systemPrompt += "\n" + cs.SystemPrompt()
	// allow plugins to add additional context to the system message.
	// this functionality may be better suited elsewhere depending on if it's useful to call "AddSystemMessage" more than once within a conversation
	// or if we choose to register plugins after the system message has been added.
	// for now, we'll just keep it here
	for _, plugin := range sesh.Plugins() {
		if prompt, ok := plugin.(openai.PluginPrompt); ok {
			systemPrompt += "\n" + prompt.SystemPrompt()
		}
	}

	sesh.AddSystemMessage(systemPrompt)

	return sesh, nil
}

// language doc fetching plugin
