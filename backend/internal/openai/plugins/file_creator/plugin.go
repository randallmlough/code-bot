package filecreator

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/randallmlough/code-bot/internal/file"
	"github.com/randallmlough/code-bot/internal/openai"
)

const (
	pluginName  = "file-creator"
	description = "it creates a file on the local machine"
)

type FileCreator struct{}

func New() openai.Plugin {
	return &FileCreator{}
}

func (p *FileCreator) Name() string {
	return pluginName
}
func (p *FileCreator) Description() string {
	return description
}

type Args struct {
	Contents string `json:"contents"`
	Filepath string `json:"filepath"`
}

func (p *FileCreator) Parameters() openai.JSONSchemaDefinition {
	return openai.JSONSchemaDefinition{
		Type: openai.JSONSchemaTypeObject,
		Properties: map[string]openai.JSONSchemaDefinition{
			"contents": {
				Type:        openai.JSONSchemaTypeString,
				Description: "result of the conversion",
			},
			"filepath": {
				Type:        openai.JSONSchemaTypeString,
				Description: "destination of the file",
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

func (p *FileCreator) Execute(fnArgs openai.JsonObject) (openai.JsonObject, error) {
	args := Args{}
	if err := mapstructure.Decode(fnArgs, &args); err != nil {
		return nil, err
	}

	fmt.Println("FILEPATH:", args.Filepath)
	err := file.CreateFile(args.Contents, "./testing.go")
	if err != nil {
		return openai.JsonObject{
			"error": "failed to create a file",
		}, err
	}
	return openai.JsonObject{
		"status": "file created",
	}, nil
}
