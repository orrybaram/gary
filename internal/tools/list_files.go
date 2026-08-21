package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ListFiles walks a directory and returns its entries as a JSON array.
var ListFiles = Definition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, list files in the current directory.",
	InputSchema: schemaFor[listFilesInput](),
	Function:    listFiles,
}

type listFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

func listFiles(input json.RawMessage) (string, error) {
	var in listFilesInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	dir := "."
	if in.Path != "" {
		dir = in.Path
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
