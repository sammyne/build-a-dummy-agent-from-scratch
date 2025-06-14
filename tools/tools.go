package tools

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

type Definition struct {
	Name        string
	Description string
	// InputSchema now map[string]any to match OpenAI's parameter schema format
	InputSchema map[string]any
	Function    func(input json.RawMessage) (string, error) // Input is JSON string from OpenAI args
}

// GenerateSchema adapted to return map[string]any
func GenerateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties:  false,
		DoNotReference:             true, // Keep definitions inline for OpenAI
		RequiredFromJSONSchemaTags: true, // Respect `jsonschema:"required"`
	}
	var v T
	schema := reflector.Reflect(v)

	// Convert the jsonschema.Schema to map[string]any expected by OpenAI
	// This is a simplification; a full conversion might be more complex
	schemaBytes, _ := json.Marshal(schema)
	var schemaMap map[string]any
	_ = json.Unmarshal(schemaBytes, &schemaMap)

	// OpenAI expects parameters schema directly, remove unnecessary outer layers if present
	if props, ok := schemaMap["properties"]; ok {
		schemaMap["properties"] = props
	}
	if req, ok := schemaMap["required"]; ok {
		schemaMap["required"] = req
	}
	schemaMap["type"] = "object" // Ensure root type is object

	return schemaMap
}
