package tools

import (
	"encoding/json"
	"os"
)

// ReadFile returns the contents of a single file.
var ReadFile = Definition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: schemaFor[readFileInput](),
	Function:    readFile,
}

type readFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

func readFile(input json.RawMessage) (string, error) {
	var in readFileInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	content, err := os.ReadFile(in.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
