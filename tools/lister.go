package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ListFiles Tool Definition
type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

var ListFilesDefinition = Definition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory. Returns a JSON array of strings, directories have a trailing slash.",
	InputSchema: GenerateSchema[ListFilesInput](),
	Function:    ListFiles, // Function implementation remains the same
}

// ListFiles function implementation unchanged
func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	if len(input) > 0 && string(input) != "null" {
		err := json.Unmarshal(input, &listFilesInput)
		if err != nil {
			return "", fmt.Errorf("failed to parse input for list_files: %w. Input was: %s", err, string(input))
		}
	}
	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}
	var files []string
	err := filepath.WalkDir(dir, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, currentPath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", currentPath, err)
		}
		if relPath == "." {
			return nil
		}
		if d.IsDir() {
			files = append(files, relPath+"/")
		} else {
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error listing files in '%s': %w", dir, err)
	}
	result, err := json.Marshal(files)
	if err != nil {
		return "", fmt.Errorf("failed to marshal file list to JSON: %w", err)
	}
	return string(result), nil
}