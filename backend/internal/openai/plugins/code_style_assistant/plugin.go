package codestyleassistant

import (
	"embed"

	"github.com/randallmlough/code-bot/internal/file"
)

//go:embed style-guide.md
var styleGuide embed.FS

const (
	pluginName  = "code-style-assistant"
	description = "it will tell you if the user has a preference on how they like to code in a language with examples."
)

type CodeStyleAssistant struct {
	styleGuide string
}

func New() (*CodeStyleAssistant, error) {
	guide, err := file.ReadEmbededFile(styleGuide, "style-guide.md")
	if err != nil {
		return nil, err
	}
	return &CodeStyleAssistant{styleGuide: guide}, nil
}

func (p *CodeStyleAssistant) SystemPrompt() string {
	return p.styleGuide
}
