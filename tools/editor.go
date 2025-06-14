package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// EditFile Tool Definition
type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file" jsonschema:"required"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for. If empty and file doesn't exist, creates the file with new_str as content. If not empty, MUST match exactly (limitation)."`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with, or the initial content if creating a new file." jsonschema:"required"`
}

var EditFileDefinition = Definition{
	Name:        "edit_file",
	Description: `Make edits to a text file. Replaces ALL occurrences of 'old_str' with 'new_str'. If 'old_str' is empty and the file doesn't exist, it creates it with 'new_str'.`,
	InputSchema: GenerateSchema[EditFileInput](),
	Function:    EditFile, // Function implementation remains the same
}

// EditFile and createNewFile function implementations unchanged
func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}
	}
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file '%s': %w", filePath, err)
	}
	return fmt.Sprintf("Successfully created file %s", filePath), nil
}
func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse input for edit_file: %w. Input was: %s", err, string(input))
	}
	if editFileInput.Path == "" {
		return "", fmt.Errorf("invalid input: 'path' cannot be empty")
	}
	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", fmt.Errorf("error reading file '%s': %w", editFileInput.Path, err)
	}
	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)
	if oldContent == newContent && editFileInput.OldStr != "" {
		if oldContent == editFileInput.NewStr {
			return "OK (no change needed, content already matched)", nil
		}
		return "", fmt.Errorf("old_str '%s' not found in file '%s'", editFileInput.OldStr, editFileInput.Path)
	}
	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing file '%s': %w", editFileInput.Path, err)
	}
	return "OK", nil
}