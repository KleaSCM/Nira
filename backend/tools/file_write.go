package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileWriteTool struct {
	AllowedPaths []string
}

func NewFileWriteTool(allowedPaths []string) *FileWriteTool {
	return &FileWriteTool{
		AllowedPaths: allowedPaths,
	}
}

func (t *FileWriteTool) Name() string {
	return "write_file"
}

func (t *FileWriteTool) Description() string {
	return "Writes text content to a file. Arguments: 'path' (string), 'content' (string)."
}

func (t *FileWriteTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.Name(),
		"description": t.Description(),
		"parameters": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to write to",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "The text content to write",
				},
			},
			"required": []string{"path", "content"},
		},
	}
}

func (t *FileWriteTool) Execute(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path argument is required and must be a string")
	}
	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content argument is required and must be a string")
	}

	if !t.isPathAllowed(path) {
		return nil, fmt.Errorf("path '%s' is not in allowed directories", path)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote to %s", path), nil
}

func (t *FileWriteTool) isPathAllowed(path string) bool {
	if len(t.AllowedPaths) == 0 {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	for _, allowed := range t.AllowedPaths {
		absAllowed, err := filepath.Abs(allowed)
		if err != nil {
			continue
		}

		rel, err := filepath.Rel(absAllowed, absPath)
		if err != nil {
			continue
		}

		if rel != ".." && !filepath.IsAbs(rel) {
			return true
		}
	}

	return false
}
