package openai

import gogpt "github.com/sashabaranov/go-openai"

type FunctionDefinition = gogpt.FunctionDefinition

type JSONSchemaType = gogpt.JSONSchemaType

const (
	JSONSchemaTypeObject  JSONSchemaType = "object"
	JSONSchemaTypeNumber  JSONSchemaType = "number"
	JSONSchemaTypeString  JSONSchemaType = "string"
	JSONSchemaTypeArray   JSONSchemaType = "array"
	JSONSchemaTypeNull    JSONSchemaType = "null"
	JSONSchemaTypeBoolean JSONSchemaType = "boolean"
)

type JSONSchemaDefinition = gogpt.JSONSchemaDefinition
