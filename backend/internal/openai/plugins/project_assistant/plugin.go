package projectassistant

import (
	"embed"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/randallmlough/code-bot/internal/file"
	"github.com/randallmlough/code-bot/internal/openai"
)

//go:embed system_prompt.md
var content embed.FS

const (
	pluginName  = "project_assistant"
	description = "it manages files on the local machine"
)

type ProjectAssistant struct {
	projectTree  files
	systemPrompt string
}

func New() *ProjectAssistant {
	ff, err := GetFileStructure(".")
	if err != nil {
		panic(fmt.Errorf("failed to walk the project structure: %w", err))
	}
	prompt, err := generateSystemPrompt()
	if err != nil {
		panic(err)
	}
	return &ProjectAssistant{
		projectTree:  ff,
		systemPrompt: prompt,
	}
}

func (p *ProjectAssistant) Name() string {
	return pluginName
}

func (p *ProjectAssistant) Description() string {
	return description
}

func (p *ProjectAssistant) UserMessage() string {
	prompt := "\nCurrent project structure: ###\n"
	prompt += p.projectTree.MarkdownTable()
	prompt += "\n###\n"
	return prompt
}

func generateSystemPrompt() (string, error) {
	prompt, err := file.ReadEmbededFile(content, "system_prompt.md")
	if err != nil {
		return "", fmt.Errorf("%s plugin failed to read the system prompt: %w", pluginName, err)
	}
	return prompt, nil
}

func (p *ProjectAssistant) SystemPrompt() string {
	return p.systemPrompt
}

type Args struct {
	Body      string `json:"body"`
	Filepath  string `json:"filepath"`
	Operation string `json:"operation"`
}

func (p *ProjectAssistant) Parameters() openai.JSONSchemaDefinition {
	return openai.JSONSchemaDefinition{
		Type: openai.JSONSchemaTypeObject,
		Properties: map[string]openai.JSONSchemaDefinition{
			"body": {
				Type:        openai.JSONSchemaTypeString,
				Description: "The contents of the request.",
			},
			"filepath": {
				Type:        openai.JSONSchemaTypeString,
				Description: "Where the file can be found or created.",
			},
			"operation": {
				Type:        openai.JSONSchemaTypeString,
				Description: "What to do with the file. Defaults to creating a file unless it already exists in which case it will update the existing file.",
				Enum:        []string{"create", "update"},
			},
			// "words": {
			// 	Type:        openai.JSONSchemaTypeArray,
			// 	Description: "list of words in sentence",
			// 	Items: &openai.JSONSchemaDefinition{
			// 		Type: openai.JSONSchemaTypeString,
			// 	},
			// },
			// "enumTest": {
			// 	Type: openai.JSONSchemaTypeString,
			// 	Enum: []string{"hello", "world"},
			// },
		},
		Required: []string{"contents", "filepath"},
	}
}

func (p *ProjectAssistant) Execute(fnArgs openai.JsonObject) (openai.JsonObject, error) {
	args := Args{}
	if err := mapstructure.Decode(fnArgs, &args); err != nil {
		return nil, err
	}

	fmt.Println("FILEPATH:", args.Filepath)
	switch args.Operation {
	case "create":
		err := createFile(args.Body, args.Filepath)
		if err != nil {
			return openai.JsonObject{
				"error": "failed to create a file",
			}, err
		}
	case "update":
		err := updateFile(args.Body, args.Filepath)
		if err != nil {
			return openai.JsonObject{
				"error": "failed to update file",
			}, err
		}
	}

	return openai.JsonObject{
		"status": "file created",
	}, nil
}
