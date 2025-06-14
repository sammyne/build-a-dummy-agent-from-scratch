package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadFile Tool Definition
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory." jsonschema:"required"`
}

var ReadFileDefinition = Definition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: GenerateSchema[ReadFileInput](),
	Function:    ReadFile, // Function implementation remains the same
}

func ReadFile(input json.RawMessage) (string, error) { // Implementation unchanged
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse input for read_file: %w. Input was: %s", err, string(input))
	}
	if readFileInput.Path == "" {
		return "", fmt.Errorf("missing required parameter 'path' for read_file")
	}
	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %w", readFileInput.Path, err)
	}
	return string(content), nil
}
