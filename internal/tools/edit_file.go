package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EditFile replaces an exact substring in a file, creating the file when
// old_str is empty and the path does not exist.
var EditFile = Definition{
	Name: "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str'
MUST be different from each other. If the file specified with path doesn't
exist, it will be created.`,
	InputSchema: schemaFor[editFileInput](),
	Function:    editFile,
}

type editFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only match once"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str"`
}

func editFile(input json.RawMessage) (string, error) {
	var in editFileInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	if in.Path == "" || in.OldStr == in.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(in.Path)
	if err != nil {
		if os.IsNotExist(err) && in.OldStr == "" {
			return createNewFile(in.Path, in.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.ReplaceAll(oldContent, in.OldStr, in.NewStr)

	if oldContent == newContent && in.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	if err := os.WriteFile(in.Path, []byte(newContent), 0644); err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	if dir := filepath.Dir(filePath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}
